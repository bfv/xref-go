 
 /*------------------------------------------------------------------------
    File        : OrderBT
    Purpose     : 
    Syntax      : 
    Description : 
    Author(s)   : bronco
    Created     : Thu Jul 12 16:34:21 CEST 2018
    Notes       : 
  ----------------------------------------------------------------------*/

using Progress.Lang.*.

block-level on error undo, throw.

class oo.OrderBT: 
  
  define public property currentOrder as integer no-undo get. set.
  
  constructor public OrderBT():
      
  end.
  
  constructor public OrderBT(orderNumber as integer):
    
  end method.
    
end class.