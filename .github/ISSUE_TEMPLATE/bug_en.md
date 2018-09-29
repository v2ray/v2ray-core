---
name: Bug report
about: "Create a bug report to help us improve"

---

Please read the [instruction](https://github.com/v2ray/v2ray-core/blob/master/.github/SUPPORT.md) and answer the following questions before submitting your issue. Thank you.
Please answer all the questions with enough information. All issues not following this template will be closed immediately.
If you are not sure if your question is truely a bug in V2Ray, please discuss it [here](https://github.com/v2ray/discussion/issues) first.

1) What version of V2Ray are you using (If you deploy different version on server and client, please explicitly point out)?

2) What's your scenario of using V2Ray? E.g., Watching YouTube videos in Chrome via Socks/VMess proxy.

3) What did you see? (Please describe in detail, such as timeout, fake TLS certificate etc)

4) What's your expectation?

5) Please attach your configuration file (**Mask IP addresses before submit this issue**).

Server configuration:

```javascript
    // Please attach your server configuration here.
```

Client configuration:

```javascript
    // Please attach your client configuration here.
```

6) Please attach error logs, especially the bottom lines if the file is large. Error log file is usually at `/var/log/v2ray/error.log` on Linux.

Server error log:

```javascript
    // Please attach your server error log here.
```

Client error log:

```javascript
    // Please attach your client error log here.
```

7) Please attach access log. Access log is usually at '/var/log/v2ray/access.log' on Linux.

```javascript
    // Please attach your server access log here.
```

Please review your issue before submitting.

8) Other configurations (such as Nginx) and logs.

9) If V2Ray doesn't run, please attach output from `--test`.

The command is usually `/usr/bin/v2ray/v2ray --test --config /etc/v2ray/config.json`, but may vary according to your scenario.

10) If V2Ray service doesn't run, please attach journal log.

Usual command is `journalctl -u v2ray`.
