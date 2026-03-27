 
 /*------------------------------------------------------------------------
    File        : Employee
    Purpose     : 
    Syntax      : 
    Description : 
    Author(s)   : bronco
    Created     : Mon Jul 02 12:01:31 CEST 2018
    Notes       : 
  ----------------------------------------------------------------------*/

using oo.person.

block-level on error undo, throw.

class oo.Employee inherits person: 
  
  define public property EmployeeNumber as integer no-undo get. set.
  
end class.
