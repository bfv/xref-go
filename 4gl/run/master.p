
define variable procje as character no-undo.
define variable hproc as handle no-undo.

procje = 'xref/oo/proc1.p'.

run xref/oo/proc1.p.
run xref/oo/proc2.p.

run value(procje).

run xref/oo/persistent.p persistent set hproc.
run test1 in hproc ('hi!').
