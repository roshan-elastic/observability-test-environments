# Sudo with Touch ID

In order to use Touch ID with sudo, you need to create the `/etc/pam.d/sudo_local`,
and add the following line to it:

```
auth       sufficient     pam_tid.so
```

You can make this change by running the following command:

```
echo "auth       sufficient     pam_tid.so" | sudo tee -a /etc/pam.d/sudo_local
```

After this change any `sudo` command will request Touch ID authentication.
