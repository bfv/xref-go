 
 /*------------------------------------------------------------------------
    File        : CustomerBE
    Purpose     : 
    Syntax      : 
    Description : 
    Author(s)   : bronco
    Created     : Thu Jul 12 16:28:15 CEST 2018
    Notes       : 
  ----------------------------------------------------------------------*/

using Progress.Lang.*.

block-level on error undo, throw.

class oo.CustomerBE: 
  
  method public void GetByName(custName as character):
    define buffer b-customer for customer.
    for each b-customer where b-customer.Name begins custName no-lock:
      // bla
    end.
  end method.
  
  method public void GetNoIndex():
    for each customer where Customer.Discount > 10 table-scan no-lock:
      // no index
    end.
  end method.
  
end class.