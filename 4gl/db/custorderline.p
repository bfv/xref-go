
define temp-table ttsports no-undo
  field name as character
  field ordernum as integer
  field itemnum as integer
  .
  
for each customer no-lock,
    each order where order.custnum = customer.custnum no-lock,
    each orderline where orderline.ordernum = order.ordernum no-lock:
  
  create ttsports.
  assign 
    ttsports.name = customer.name 
    ttsports.ordernum = order.ordernum 
    ttsports.itemnum = OrderLine.Itemnum
    .
        
end.
