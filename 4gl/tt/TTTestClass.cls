 
 /*------------------------------------------------------------------------
    File        : TTTestClass
    Purpose     : 
    Syntax      : 
    Description : 
    Author(s)   : bronco
    Created     : Tue Jul 17 16:47:32 CEST 2018
    Notes       : 
  ----------------------------------------------------------------------*/

using Progress.Lang.*.

block-level on error undo, throw.

class tt.TTTestClass: 

  {tt/tttest.i &accessor="protected"}
  
  method public void updateAccess():
    find first tttest where tttest.ctest = "hoppa".
    assign tttest.itest = 5.    
  end method.
  
end class.