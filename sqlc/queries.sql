-- name: ListFunctions
SELECT id, file_id, name, signature FROM functions WHERE repo_id = $1;

-- name: ListStructs
SELECT id, file_id, name FROM structs WHERE repo_id = $1;

-- name: ListEdges
SELECT from_id, to_id, type FROM edges WHERE repo_id = $1;

-- name: ListHotFiles
SELECT file_id, count(*) as changes FROM commits_files WHERE repo_id = $1 GROUP BY file_id ORDER BY changes DESC LIMIT 10;
