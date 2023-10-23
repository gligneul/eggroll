---
title: FAQ
---

# EggRoll FAQ

## Why Go?

The Cartesi Machine, a high-performance RISC-V emulator, excels in its domain.
Yet, no emulator can truly rival bare metal performance.
Utilizing scripting languages like JavaScript or Python within the Cartesi Machine results in double emulation—first for the Cartesi Machine and then for the language interpreter.
To address this, one should use languages that compile directly into RISC-V machine code, like Go, Rust, or C++.
We selected Go for the EggRoll framework due to its exceptional performance and ease of adoption.
