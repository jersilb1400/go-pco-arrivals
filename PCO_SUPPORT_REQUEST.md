# PCO Support Request: Check-Ins API Access

## Subject
Request for Check-Ins API Access for OAuth Application

## Application Details
- **Application Name**: PCO Arrivals Billboard
- **Client ID**: `8a7f86cc34fb86744d5c36a0d5a4af8a1a5447aa030f94653037c5ba909aced5`
- **Application Type**: OAuth 2.0 Application
- **Use Case**: Church check-in display system for real-time arrival notifications

## Issue Description
Our OAuth application is requesting the `check_ins` scope but receiving 401 Unauthorized errors when attempting to access the Check-Ins API endpoints.

### Error Details
```
HTTP Status: 401 Unauthorized
Error Code: bad_scope
Error Title: Request outside authenticated scope
Error Detail: The API credentials do not have access to the application check-ins
```

### Requested Scopes
- `people` ✅ (working)
- `check_ins` ❌ (failing)

## Steps Taken
1. ✅ Verified OAuth application is properly configured
2. ✅ Confirmed `check_ins` scope is requested in authorization URL
3. ✅ Tested with valid access tokens
4. ✅ Checked application permissions in PCO Developer Dashboard

## Request
Please grant our application access to the Check-Ins API so we can:
- Retrieve check-in events for specific dates
- Access check-in locations and times
- Build a real-time arrival display system for our church

## Contact Information
- **Organization**: Grace Fellowship
- **Use Case**: Church check-in display system
- **Technical Contact**: [Your Name/Email]

---

**Note**: This is for internal church use to display real-time check-ins on arrival screens. 