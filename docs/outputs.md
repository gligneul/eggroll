---
title: Generating Outputs
---

Generating Outputs
=

### Voucher
- **Purpose**: A voucher is a combination of a target address and a payload in bytes. It is used by the off-chain machine to respond and interact with L1 smart contracts. Upon execution, a voucher sends a message to the target address with the payload as a parameter.
- **Example**:
  ```go
  func (c *Contract) AdvanceFunction(env eggroll.Env, value string) error {
     env.Voucher(EncodeEchoResponse(value))
      return nil
  }
  ```

### Notice
- **Purpose**: A notice is an arbitrary payload in bytes that is submitted by the off-chain machine for informational purposes. Similarly to vouchers, when the epoch containing a notice is finalized a proof will be produced so that the validity of its content can be verified on-chain by any interested party.
- **Example**:
  ```go
  func (c *Contract) AdvanceFunction(env eggroll.Env, value string) error {
      env.Notice(EncodeEchoResponse(value))
      return nil
  }
  ```

### Report
- **Purpose**: A report is an application log or a piece of diagnostic information. Like a notice, it is represented by an arbitrary payload in bytes. However, a report is never associated with a proof and is thus not suitable for trustless interactions such as on-chain processing or convincing independent third parties of dApp outcomes. Reports are commonly used to indicate processing errors or to retrieve application information for display.
- **Example**:
  ```go
  func (c *Contract) AdvanceEcho(env eggroll.Env, value string) error {
      env.Report(EncodeEchoResponse(value))
      return nil
  }
  ```

