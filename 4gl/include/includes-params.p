
{include/base.i}
{include/dsfamily.i &test="testing12"}

create ttparent.
assign 
  ttparent.parent_id = 1
  ttparent.parentname = 'parent1'
  .
  
create ttchild.
assign
  ttchild.child_id = 1
  ttchild.parent_id = 1
  ttchild.childname = 'child1'
  .
  
create ttchild.
assign
  ttchild.child_id = 2
  ttchild.parent_id = 1
  ttchild.childname = 'child2'
  .
  
