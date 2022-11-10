package assist

const (
	AdClassDomainDNS          = "domainDNS"
	AdClassContainer          = "container"
	AdClassOrganizationalUnit = "organizationalUnit"
	AdClassComputer           = "computer"
	AdClassContact            = "contact"
	AdClassGroup              = "group"
	AdClassUser               = "user"
)

const (
	AdCategoryOrganizationalUnit = "OrganizationalUnit"
	AdCategoryContainer          = "Container"
	AdCategoryComputer           = "Computer"
	AdCategoryPerson             = "Person"
	AdCategoryGroup              = "Group"
)

const (
	AdUserAccountEnable  = "66048"
	AdUserAccountDisable = "66050"
)

const (
	AdScript                             = 0x00000001 // The logon script is executed.
	AdAccountDisable                     = 0x00000002 // The user account is disabled.
	AdHomedirRequired                    = 0x00000008 // The home directory is required.
	AdLockout                            = 0x00000010 // The account is currently locked out.
	AdPasswdNotReqD                      = 0x00000020 // No password is required.
	AdPasswdCantChange                   = 0x00000040 // The user cannot change the password.
	AdEncryptedTextPasswordAllowed       = 0x00000080 // The user can send an encrypted password.
	AdTempDuplicateAccount               = 0x00000100 // This is an account for users whose primary account is in another domain.
	AdNormalAccount                      = 0x00000200 // This is a default account type that represents a typical user.
	AdInterDomainTrustAccount            = 0x00000800 // This is a permit to trust account for a system domain that trusts other domains.
	AdWorkstationTrustAccount            = 0x00001000 // This is a computer account for a computer that is a member of this domain.
	AdServerTrustAccount                 = 0x00002000 // This is a computer account for a system backup domain controller that is a member of this domain.
	AdUnused1                            = 0x00004000 // Not used.
	AdUnused2                            = 0x00008000 // Not used.
	AdDontExpirePasswd                   = 0x00010000 // The password for this account will never expire.
	AdMnsLogonAccount                    = 0x00020000 // This is an MNS logon account.
	AdSmartCardRequired                  = 0x00040000 // The user must log on using a smart card.
	AdTrustedForDelegation               = 0x00080000 // The service account (user or computer account), under which a service runs, is trusted for Kerberos delegation.
	AdNotDelegated                       = 0x00100000 // The security context of the user will not be delegated to a service even if the service account is set as trusted for Kerberos delegation.
	AdUseDesKeyOnly                      = 0x00200000 // Restrict this principal to use only Data Encryption Standard (DES) encryption types for keys.
	AdDontRequirePreAuth                 = 0x00400000 // This account does not require Kerberos pre-authentication for logon.
	AdPasswordExpired                    = 0x00800000 // The user password has expired. This flag is created by the system using data from the Pwd-Last-Set attribute and the domain policy.
	AdTrustedToAuthenticateForDelegation = 0x01000000 // The account is enabled for delegation.
	AdUseAesKeys                         = 0x08000000
)
