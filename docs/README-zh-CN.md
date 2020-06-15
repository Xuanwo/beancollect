# beancollect

beancollect 帮助你收集豆子 (bean) 以便于清点他们。

## 使用方法

```bash
beancollect -s wechat collect/wechat.csv
```

将会输出：

```
2019-06-20 ! "北京麦当劳食品有限公司" "北京麦当劳食品有限公司"
    Assets:Deposit:WeChat -20.5 CNY
    Expenses:Intake:FastFood
2019-06-23 ! "街电科技" "街电充电宝"
    Assets:Deposit:WeChat -4 CNY
```

## 设置

推荐的 beancount 项目结构：

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

beancollect 需要的所有文件都在 `collect` 目录下。

## 配置

beancollect 支持账户映射和规则执行。

beancollect 会有一个全局的配置文件 `global.yaml`，然后每种格式都会一个独立的配置文件 `schema.yaml`，`global.yaml` 中的条目将会被 `schema.yaml` 覆盖。

比如：

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

这个配置会用来做账户的映射。

### 规则

- add_accounts: 如果条件满足，就在事务中增加一个账户

## 账单

### 微信

`我` -> `支付` -> `钱包` -> `账单` -> `...` in right up corner -> `导出账单`

账单将会发送到你指定的邮箱。

### 支付宝

访问 `https://www.alipay.com/` 下载帐单

### 招商银行信用卡

启用电子帐单，并将邮件中的 HTML 下载下来