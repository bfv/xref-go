 
 /*------------------------------------------------------------------------
    File        : DeliverAddress
    Purpose     : 
    Syntax      : 
    Description : 
    Author(s)   : bronco
    Created     : Tue Jul 03 08:43:12 CEST 2018
    Notes       : 
  ----------------------------------------------------------------------*/

using oo.Address.

block-level on error undo, throw.

class oo.DeliverAddress inherits Address final: 
  
  define public property Floor as integer no-undo get. set.
  define temp-table tttest no-undo
    field itest as integer
    field ctest as character
    index itest as primary unique itest
    .
  
  constructor DeliverAddress():
  end.
  
  constructor DeliverAddress(address as Address):
    
  end.
  
  constructor protected DeliverAddress(countAddress as integer):
    
  end.
  
  constructor static DeliverAddress(): 
    
  end.
  
  method public final character MultipleParameters(cIn as character, iIn as integer, output cOut as Address, input-output testIO as character):
    cOut = new Address().
    return cIn + string(iIn).
  end.
  
  method public void BufferParameter(buffer b-tttest for tttest):
  
  end method.
  
  method protected void TableParameter(output table tttest):
    
  end method.
  
end class.