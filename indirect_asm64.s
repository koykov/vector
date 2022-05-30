TEXT    ·indirectVector1(SB), $0-16
PCDATA  $0, $-2
MOVQ    addr+8(SP), AX
PCDATA  $0, $-1
MOVQ    AX, r1+16(SP)
RET

TEXT    ·indirectNode1(SB), $0-16
PCDATA  $0, $-2
MOVQ    addr+8(SP), AX
PCDATA  $0, $-1
MOVQ    AX, r1+16(SP)
RET
