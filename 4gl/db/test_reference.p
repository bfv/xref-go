
define temp-table ttcustomer no-undo
  field custnum like Customer.CustNum
  .
  
find first Customer no-lock.
create ttcustomer.
assign ttcustomer.custnum = Customer.CustNum.
  