 
 /*------------------------------------------------------------------------
    File        : MultipleImplements
    Purpose     : 
    Syntax      : 
    Description : 
    Author(s)   : bronco
    Created     : Wed Jul 18 08:15:42 CEST 2018
    Notes       : 
  ----------------------------------------------------------------------*/

using Progress.Lang.*.
using oo.IEmpty.
using oo.IDisposable.

block-level on error undo, throw.

class oo.MultipleImplements implements IEmpty, IDisposable: 
  
  method public void Dispose():
    
  end method.
  
end class.