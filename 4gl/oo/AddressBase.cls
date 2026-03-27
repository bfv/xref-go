 
 /*------------------------------------------------------------------------
    File        : AddressBase
    Purpose     : 
    Syntax      : 
    Description : 
    Author(s)   : bronco
    Created     : Mon Jul 02 16:08:13 CEST 2018
    Notes       : 
  ----------------------------------------------------------------------*/

using Progress.Lang.*.
using oo.IBase.
using oo.Address.

block-level on error undo, throw.

class oo.AddressBase implements IBase abstract: 
  
  method public abstract void MoveTo(adress as Address).
  method public character Identify(name as character). end method.
  
end class.