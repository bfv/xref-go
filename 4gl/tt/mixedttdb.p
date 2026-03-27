
{tt/ttsports.i}

create ttsports.
assign
  ttsports.name = 'abc'
  ttsports.ordernum = 1
  ttsports.itemnum = 2
  .

run testMixed.

procedure testMixed:
  
  find first ttsports.
  for each order where order.ordernum = ttsports.ordernum no-lock:
    display ttsports.name Order.OrderStatus.
  end.  
  
end procedure.