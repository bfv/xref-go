 
 /*------------------------------------------------------------------------
    File        : TTAccess
    Purpose     : 
    Syntax      : 
    Description : 
    Author(s)   : bronco
    Created     : Tue Jul 17 16:02:47 CEST 2018
    Notes       : 
  ----------------------------------------------------------------------*/

using Progress.Lang.*.

block-level on error undo, throw.

class tt.TTAccess: 

  {tt/ttsports.i &accessor="protected"}
  
  method public void updateAccess():
    find first ttsports where ttsports.ordernum = 1.
    assign ttsports.itemnum = 5.    
  end method.
  
end class.