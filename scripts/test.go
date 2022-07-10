package scripts

var Test = `
set @code := ?;
set @name := ?;

set @safe_id := REPLACE(UUID(),'-','');

#insert
insert INTO test_table SET 
ID=@safe_id, 
CODE=@code,
NAME=@name;

#select
SELECT * FROM test_table WHERE ID=@safe_id;
`
