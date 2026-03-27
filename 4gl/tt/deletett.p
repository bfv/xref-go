
{tt/ttsports.i}

for each ttsports:
  delete ttsports.
end.

// just a create with no find's doesn't show in the XREF
find first ttsports.
delete ttsports.

procedure publicDelete:
  for each ttsports:
    delete ttsports.
  end.
end procedure.
