
{tt/ttsports.i}

// just a create with no find's doesn't show up in the XREF
create ttsports.
assign
  ttsports.name = 'abc'
  ttsports.ordernum = 1
  ttsports.itemnum = 2
  .

procedure publicCreate:
  
  // just a create with no find's doesn't show up in the XREF
  create ttsports.
  assign
    ttsports.name = 'abc'
    ttsports.ordernum = 1
    ttsports.itemnum = 2
    .

end procedure.
