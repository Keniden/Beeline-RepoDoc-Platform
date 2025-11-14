package git

import (
    "context"
    "fmt"

    gogit "github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing/object"
)

type Report struct {
    Authors map[string]int
    Commits int
    FileCoupling map[string]map[string]int
    CommitRanges map[string][]string
}

type Analyzer struct{
}

func NewAnalyzer() *Analyzer {
    return &Analyzer{}
}

func (a *Analyzer) Run(ctx context.Context, repoPath string) (*Report, error) {
    repo, err := gogit.PlainOpen(repoPath)
    if err != nil {
        return nil, fmt.Errorf("open repo: %w", err)
    }
    ref, err := repo.Head()
    if err != nil {
        return nil, fmt.Errorf("repo head: %w", err)
    }
    iter, err := repo.Log(&gogit.LogOptions{From: ref.Hash()})
    if err != nil {
        return nil, fmt.Errorf("log: %w", err)
    }
    report := &Report{
        Authors: make(map[string]int),
        FileCoupling: make(map[string]map[string]int),
        CommitRanges: make(map[string][]string),
    }
    err = iter.ForEach(func(c *object.Commit) error {
        report.Commits++
        report.Authors[c.Author.Email]++
        files := map[string]struct{}{}
        stats, err := c.Stats()
        if err != nil {
            return err
        }
        for _, stat := range stats {
            files[stat.Name] = struct{}{}
        }
        for a := range files {
            if report.CommitRanges[a] == nil {
                report.CommitRanges[a] = []string{}
            }
            report.CommitRanges[a] = append(report.CommitRanges[a], c.Hash.String())
            if report.FileCoupling[a] == nil {
                report.FileCoupling[a] = make(map[string]int)
            }
            for b := range files {
                if a == b {
                    continue
                }
                report.FileCoupling[a][b]++
            }
        }
        return nil
    })
    if err != nil {
        return nil, fmt.Errorf("iterate commits: %w", err)
    }
    return report, nil
}
