
run test1('before').

procedure test1:
  define input parameter cIn as character no-undo.
end procedure.

procedure testPrivate private:
  
end procedure.

run test1('after').
