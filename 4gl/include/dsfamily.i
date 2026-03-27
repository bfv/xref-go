
{include/ttparent.i {&*}}
{include/ttchild.i {&*}}

define {&accessor} dataset dsfamily 
  for ttparent, ttchild
  data-relation parent_child 
    for ttparent, ttchild
    relation-fields (parent_id, parent_id)
    . 