 
 /*------------------------------------------------------------------------
    File        : Person
    Purpose     : 
    Syntax      : 
    Description : 
    Author(s)   : bronco
    Created     : Mon Jul 02 12:00:41 CEST 2018
    Notes       : 
  ----------------------------------------------------------------------*/

using oo.Address.

block-level on error undo, throw.

class oo.Person: 
  
  define public property FirstName as character no-undo get. set.
  define public property LastName as character no-undo get. set.
  define public property BirthDate as date no-undo get. set.
  
  define public property HomeAddress as Address no-undo get. set.
  
  constructor public Person():
    this-object:HomeAddress = new Address().  
  end.
  
  method public override character ToString():
    return substitute("&1~n&2", this-object:FirstName + ' ' + this-object:LastName, this-object:HomeAddress:GetAdress()).  
  end method.
  
  method public character GetHomeCity():
    this-object:HomeAddress:ToString().
    return this-object:HomeAddress:City.
  end method.
  
end class.
