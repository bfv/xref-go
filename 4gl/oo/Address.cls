 
 /*------------------------------------------------------------------------
    File        : Address
    Purpose     : 
    Syntax      : 
    Description : 
    Author(s)   : bronco
    Created     : Mon Jul 02 12:07:15 CEST 2018
    Notes       : 
  ----------------------------------------------------------------------*/

using Progress.Lang.*.
using oo.AddressBase from propath.
using oo.Address from propath.

block-level on error undo, throw.

class oo.Address inherits AddressBase: 
  
  define public property Street as character no-undo get. set.
  define public property City as character no-undo get. set.
  
  method public character GetAdress():
    return this-object:Street + ' ' + this-object:City.
  end method.
  
  method public override character ToString():
    return super:ToString().
  end method.
  
  method public override void MoveTo(address as Address):
    
  end method.
  
end class.
