
{tt/ttsports.i}

define input-output parameter table for ttsports.

find first ttsports where ttsports.ordernum = 1.
assign ttsports.itemnum = 5.

procedure publicAssign:
  find first ttsports where ttsports.ordernum = 1.
  assign ttsports.itemnum = 5.
end procedure.