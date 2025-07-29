# PCO New OAuth App Setup Checklist

## 🎯 **Step 1: Create New PCO OAuth Application**

### **Go to PCO Developer Dashboard**
- URL: https://api.planningcenteronline.com/oauth/applications
- Sign in with your PCO account

### **Create New Application**
- Click **"New Application"**
- Fill in the following details:

| Field | Value |
|-------|-------|
| **Name** | `PCO Arrivals Billboard v2` |
| **Description** | `Church check-in display system for real-time arrival notifications` |
| **Website** | `http://localhost:3000` |
| **Callback URLs** | `http://localhost:3000/auth/callback` |

### **Configure Scopes**
**CRITICAL**: Check these scopes in the application settings:

- ✅ **People** - Required for user authentication
- ✅ **Check-Ins** - Required for events and check-in data
- ✅ **Services** - Optional, for additional service data
- ✅ **Groups** - Optional, for group management

### **Save Application**
- Click **"Create Application"**
- Note down the **Client ID** and **Client Secret**

---

## 🔧 **Step 2: Update Local Configuration**

### **Run the Update Script**
```bash
chmod +x update_pco_credentials.sh
./update_pco_credentials.sh
```

### **Or Manually Update .env File**
```bash
# Edit .env file
PCO_CLIENT_ID=your_new_client_id_here
PCO_CLIENT_SECRET=your_new_client_secret_here
PCO_SCOPES=people check_ins
```

---

## 🚀 **Step 3: Restart and Test**

### **Restart Backend Server**
```bash
# Stop current server (Ctrl+C)
# Then restart
go run main.go
```

### **Clear Existing Sessions**
```bash
# Clear browser cookies for localhost:3000
# Or use incognito/private browsing
```

### **Test the New Setup**
```bash
chmod +x test_pco_status.sh
./test_pco_status.sh
```

---

## ✅ **Step 4: Verification Checklist**

- [ ] New OAuth app created in PCO dashboard
- [ ] Check-Ins scope is enabled
- [ ] Client ID and Secret copied
- [ ] .env file updated with new credentials
- [ ] Backend server restarted
- [ ] Browser cookies cleared
- [ ] Fresh login completed
- [ ] API tests pass

---

## 🐛 **Troubleshooting**

### **If Check-Ins API Still Returns 401**
1. **Double-check scopes** in PCO dashboard
2. **Wait 5-10 minutes** for changes to propagate
3. **Contact PCO support** with new Client ID
4. **Use enhanced mock data** for development

### **If Login Fails**
1. **Verify callback URLs** match exactly
2. **Check Client ID/Secret** are correct
3. **Ensure .env file** is in project root
4. **Restart server** after credential changes

---

## 📞 **PCO Support Contact**

If you still get scope issues with the new app:

**Subject**: `Check-Ins API Access Request - New OAuth App`

**Include**:
- New Client ID: `[your_new_client_id]`
- Application Name: `PCO Arrivals Billboard v2`
- Error: `"The API credentials do not have access to the application check-ins"`
- Use Case: Church check-in display system

**Contact**: PCO Developer Support 