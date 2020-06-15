# beancollect

beancollect helps your collect beans so that you can count them.

[中文文档](docs/README-zh-CN.md)

## Usage

```bash
beancollect -s wechat collect/wechat.csv
```

Output will be like:

```
2019-06-20 ! "北京麦当劳食品有限公司" "北京麦当劳食品有限公司"
    Assets:Deposit:WeChat -20.5 CNY
    Expenses:Intake:FastFood
2019-06-23 ! "街电科技" "街电充电宝"
    Assets:Deposit:WeChat -4 CNY
```

## Setup

Recommended beancount directory structure:

```
├── account
│   ├── assets.bean
│   ├── equity.bean
│   ├── expenses.bean
│   ├── incomes.bean
│   └── liabilities.bean
├── collect
│   ├── global.yaml
│   └── wechat.yaml
├── main.bean
└── transactions
    └── 2019
        ├── 03.bean
        ├── 04.bean
        ├── 05.bean
        ├── 06.bean
        └── 07.bean
```

All file that beancollect located in `collect` directory.

## Config

beancollect supports account mapping and rules execution on beancount transactions.

beancollect will one `global.yaml` and a `schema.yaml` for every schema, and `global.yaml` will be override by `schema.yaml`.

For example:

```yaml
account:
  "招商银行(XXXX)": "Liabilities:Credit:CMB"
  "招商银行": "Assets:Deposit:CMB:CardXXXX"
  "零钱通": "Assets:Deposit:WeChat"
  "零钱": "Assets:Deposit:WeChat"

rules:
  - type: add_accounts
    condition:
      payee: "猫眼/格瓦拉生活"
    value: "Expenses:Recreation:Movie"
  - type: add_accounts
    condition:
      payee: "北京麦当劳食品有限公司"
    value: "Expenses:Intake:FastFood"
  - type: add_accounts
    condition:
      payee: "滴滴出行"
    value: "Expenses:Transport:Taxi"
  - type: add_accounts
    condition:
      payee: "摩拜单车"
    value: "Expenses:Transport:Bicycle"
```

### Schema

- wechat

### Account

Account will convert account in beancount.

### Rules

- add_accounts: If condition is matched, we will add an account into transaction.

## Billings

### WeChat

On the phone:

`Me` -> `WeChat Pay` -> `Wallet` -> `Transactions` -> `...` in right up corner -> `导出账单`

The billing will be sent to your email.

### Alipay

Visit website: `https://www.alipay.com/` to download billings

### CMBChina

Enable email billings
