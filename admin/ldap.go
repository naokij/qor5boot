package admin

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/go-ldap/ldap/v3"
	"github.com/naokij/qor5boot/models"
)

// LDAP配置环境变量
var (
	// 基本配置
	ldapEnabled      = getEnvWithDefault("LDAP_ENABLED", "false") == "true"
	ldapServer       = getEnvWithDefault("LDAP_SERVER", "")
	ldapPort         = getEnvWithDefault("LDAP_PORT", "389")
	ldapBindDN       = getEnvWithDefault("LDAP_BIND_DN", "")
	ldapBindPassword = getEnvWithDefault("LDAP_BIND_PASSWORD", "")
	ldapSearchBase   = getEnvWithDefault("LDAP_SEARCH_BASE", "")

	// 使用email作为默认搜索属性
	ldapSearchFilter = getEnvWithDefault("LDAP_SEARCH_FILTER", "(mail=%s)")

	// TLS配置
	ldapUseTLS     = getEnvWithDefault("LDAP_USE_TLS", "false") == "true"
	ldapSkipVerify = getEnvWithDefault("LDAP_SKIP_VERIFY", "false") == "true"
	ldapCertFile   = getEnvWithDefault("LDAP_CERT_FILE", "")
)

// 初始化LDAP配置，在应用启动时调用
func initLDAP() {
	// 检查LDAP配置
	if ldapEnabled {
		log.Printf("初始化LDAP配置: 启用LDAP认证")
		log.Printf("LDAP服务器: %s:%s", ldapServer, ldapPort)
		log.Printf("LDAP搜索基础: %s", ldapSearchBase)
		log.Printf("LDAP搜索过滤器: %s", ldapSearchFilter)
		log.Printf("LDAP使用TLS: %v", ldapUseTLS)

		// 验证必要配置
		if ldapServer == "" {
			log.Printf("警告: LDAP已启用，但服务器地址为空")
			ldapEnabled = false
		}

		if ldapSearchBase == "" {
			log.Printf("警告: LDAP已启用，但搜索基础为空")
		}
	} else {
		log.Printf("LDAP认证未启用")
	}

	// 设置LDAP配置到models包中
	models.SetLDAPConfig(ldapEnabled, ldapServer, authenticateWithLDAP)
}

// authenticateWithLDAP 通过LDAP认证用户
func authenticateWithLDAP(email, password string) (bool, error) {
	if !ldapEnabled || ldapServer == "" {
		log.Printf("LDAP认证未启用: enabled=%v, server=%s", ldapEnabled, ldapServer)
		return false, fmt.Errorf("LDAP认证未启用")
	}

	log.Printf("开始LDAP认证: 用户=%s, 服务器=%s:%s, 搜索基础=%s, 过滤器=%s, 使用TLS=%v",
		email, ldapServer, ldapPort, ldapSearchBase, ldapSearchFilter, ldapUseTLS)

	// 连接LDAP服务器
	log.Printf("尝试连接LDAP服务器: %s:%s", ldapServer, ldapPort)
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", ldapServer, ldapPort))
	if err != nil {
		log.Printf("LDAP连接错误: %v", err)
		return false, err
	}
	defer conn.Close()

	// 启用TLS (如果配置为使用TLS)
	if ldapUseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: ldapSkipVerify,
		}

		// 如果指定了证书文件，加载证书
		if ldapCertFile != "" {
			cert, err := os.ReadFile(ldapCertFile)
			if err != nil {
				log.Printf("读取LDAP证书文件错误: %v", err)
				return false, err
			}

			certPool := x509.NewCertPool()
			if !certPool.AppendCertsFromPEM(cert) {
				log.Printf("解析LDAP证书错误")
				return false, fmt.Errorf("无法解析LDAP证书")
			}

			tlsConfig.RootCAs = certPool
		}

		// 启动TLS
		log.Printf("开始TLS连接, 跳过验证=%v", ldapSkipVerify)
		err = conn.StartTLS(tlsConfig)
		if err != nil {
			log.Printf("LDAP StartTLS错误: %v", err)
			return false, err
		}
	}

	// 绑定服务账号（如果配置了）
	if ldapBindDN != "" && ldapBindPassword != "" {
		log.Printf("使用服务账号绑定: %s", ldapBindDN)
		err = conn.Bind(ldapBindDN, ldapBindPassword)
		if err != nil {
			log.Printf("LDAP服务账号绑定错误: %v", err)
			return false, err
		}
		log.Printf("服务账号绑定成功")
	} else {
		log.Printf("未配置服务账号，将使用匿名绑定")
	}

	// 使用email搜索用户
	searchFilter := fmt.Sprintf(ldapSearchFilter, email)
	searchReq := ldap.NewSearchRequest(
		ldapSearchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		[]string{"dn", "mail", "sAMAccountName"}, // 获取DN和其他一些属性
		nil,
	)

	log.Printf("搜索LDAP用户: 过滤器=%s, 搜索基础=%s", searchFilter, ldapSearchBase)

	result, err := conn.Search(searchReq)
	if err != nil {
		log.Printf("LDAP搜索错误: %v", err)
		return false, err
	}

	if len(result.Entries) == 0 {
		log.Printf("未找到匹配的LDAP用户")
		return false, nil
	}

	if len(result.Entries) > 1 {
		log.Printf("找到多个匹配的LDAP用户: %d个", len(result.Entries))
		for i, entry := range result.Entries {
			log.Printf("  用户 %d: DN=%s", i+1, entry.DN)
		}
		return false, nil
	}

	userDN := result.Entries[0].DN
	log.Printf("找到LDAP用户, DN: %s", userDN)

	// 输出找到的用户的所有属性，便于调试
	entry := result.Entries[0]
	log.Printf("用户属性:")
	for _, attr := range entry.Attributes {
		log.Printf("  %s: %v", attr.Name, attr.Values)
	}

	// 尝试使用用户凭据绑定
	log.Printf("尝试使用用户凭据绑定")
	err = conn.Bind(userDN, password)
	if err != nil {
		// 密码错误或其他绑定问题
		log.Printf("LDAP用户绑定失败: %v", err)
		return false, nil
	}

	// 认证成功
	log.Printf("LDAP认证成功")
	return true, nil
}
