declare namespace App {
  interface Role {
    id: number
    name: string
    code: string
    system_id: number
    system?: System
    description?: string
    is_active: boolean
    is_system_role?: boolean
  }
}

declare namespace App {
  interface ApiMeta {
    total: number
    page: number
    size: number
    pages: number
    has_next: boolean
    has_prev: boolean
    start_idx: number
    end_idx: number
  }

  interface ApiLinks {
    self: string
    next?: string
    prev?: string
  }

  interface ListData<T> {
    meta: ApiMeta
    links: ApiLinks
    items: T[]
  }

  interface IdAndName {
    id: number
    name: string
  }

  interface IdNameDesc extends IdAndName {
    description: string
  }

  interface IdNameDescSeqKeyIconCodeActive extends IdNameDesc {
    sequence: number
    key: string
    icon: string
    code: string
    is_active: boolean
  }

  interface System {
    id: number
    key: string
    icon: string
    code: string
    name: string
    description?: string | null
    is_active?: boolean | null
    path?: string | null
    sequence: number
  }

  interface Module extends IdNameDesc {
    is_active: boolean | null | undefined
    code: string
    icon: string
    system_id?: number
    system?: System
    permissions?: Permission[]
  }

  interface Model extends IdNameDesc {
    created_date?: string
  }

  interface Role {
    id: number
    name: string
    code: string
    system_id: number
    system?: System
    description?: string
    is_active: boolean
    is_system_role?: boolean
  }

  interface RoleUser {
    user_id: number
    user: User
    role_id: number
  }

  interface UserRole {
    user_id: number
    role: Role
    role_id: number
  }

  interface Permission {
    id: number
    key: string
    code: string
    name: string
    description: string
    module_id?: number
    module?: Module
    is_active?: boolean
    created_date?: string
    updated_date?: string
    action_id?: number
  }

  interface Action {
    id: number
    key: string
    code: string
    name: string
    description?: string
    is_active?: boolean | null
    created_date?: string
    updated_date?: string
  }

  interface OrganizationType extends IdNameDesc {
    code: string
  }

  interface Response<T> extends Partial<T> {
    message: string
  }

  interface Organization {
    id: number
    reg_no: string
    name: string
    short_name: string
    phone_no: string
    email: string
    logo_image_url: string
    aimag_id: number
    aimag_name: string
    sum_id: number
    sum_name: string
    bag_id: number
    bag_name: string
    address_detail: string
    country_code: string
    country_name: string
    country_name_en: string

    type_id: number
    type?: OrganizationType
    parent_id: number | null
    children?: Organization[]
    // Optional fields from extended API responses
    is_active?: boolean
    address?: string
    website?: string
    description?: string
    code?: string
  }

  interface User {
    id: number
    civil_id: number
    reg_no: string
    family_name: string
    last_name: string
    first_name: string
    gender: number
    birth_date: string
    phone_no: string
    email: string
    profile_img_url: string
  }

  interface UserDetail extends User {
    is_foreign: number
    country_code: string
    hash: string
    parent_address_id: number
    parent_address_name: string
    aimag_id: number
    aimag_code: string
    aimag_name: string
    sum_id: number
    sum_code: string
    sum_name: string
    bag_id: number
    bag_code: string
    bag_name: string
    address_detail: string
    address_type: string
    address_type_name: string
    nationality: string
    country_name: string
    country_name_en: string
  }

  interface UserProfileRes {
    is_org: boolean
    user: UserDetail
    org: Organization
  }

  interface UserOrganizationRes {
    items: Organization[]
    org: Organization
  }

  interface Scope {
    id: number
    owner_client_id: string
    key: string
    description: string
    created_date?: string
  }

  interface Notification {
    tenant: 'template.gerege.mn'
    group_id: number
    id: number
    title: string
    content: string
    user_id: number
    is_read: boolean
    idempotency_key: string
    created_date: string | null
  }

  interface AppServiceIcon extends IdAndName {
    icon: string
    link: string
    key: string
    description: string
    icon: string
    link: string
    web_link: string
    sequence: number
    is_native: boolean
    is_public: boolean
    parent_id: number
    parent?: AppServiceIcon
    children: AppServiceIcon[]
  }
  interface UserOrganization {
    birth_date: string
    first_name: string
    gender: number
    last_name: string
    org_id: number
    reg_no: string
    user_id: number
    created_date?: string
  }
  interface AppIconGroup {
    id: number
    name: string
    name_en: string
    icon: string
    type_name: string
    seq: number
    created_date: string
    updated_date: string
  }

  interface AppIcon {
    id: number
    name: string
    name_en: string
    icon: string
    icon_app: string
    icon_tablet: string
    icon_kiosk: string
    link: string
    group_id: number
    group?: AppIconGroup
    seq: number
    is_native: boolean
    is_public: boolean
    is_featured: boolean
    featured_icon: string
    is_best_selling: boolean
    feature_seq: number
    description: string
    system_code: string
    created_date: string
    updated_date: string
    is_group: boolean
    parent_id: number
    parent: AppIcon
    childs: AppIcon[]
    web_link: string
  }

  interface Menu {
    id: number
    key: string
    code: string
    name: string
    description?: string | null
    icon?: string | null
    path?: string | null
    sequence: number
    parent_id?: number | null
    permission_id?: number | null
    parent?: Menu | null
    children?: Menu[]
    system_id?: number | null
    system?: System | null
    is_active?: boolean | null
    created_date?: string | null
    updated_date?: string | null
  }

  // ============================================================
  // TPAY - Account, Card, Transaction interfaces
  // ============================================================

  /** TPay данс */
  interface TpayAccount {
    id: number
    account_no: string
    account_name: string
    balance: number
    currency: string
    is_default: boolean
    status: 'active' | 'inactive' | 'frozen'
    user_id: number
    created_date?: string
    updated_date?: string
  }

  /** Дансны хуулга */
  interface TpayStatement {
    id: number
    account_id: number
    transaction_id: string
    type: 'credit' | 'debit'
    amount: number
    balance_after: number
    description: string
    reference?: string
    counterparty_name?: string
    counterparty_account?: string
    created_date: string
  }

  /** QR код response */
  interface TpayQRCode {
    qr_string: string
    qr_image_url?: string
    amount?: number
    expires_at?: string
    account_id: number
  }

  /** Карт */
  interface TpayCard {
    id: number
    card_no: string
    card_holder_name: string
    card_type: 'visa' | 'mastercard' | 'unionpay'
    expiry_date: string
    is_verified: boolean
    is_default: boolean
    bank_name?: string
    status: 'active' | 'pending' | 'expired' | 'blocked'
    user_id: number
    created_date?: string
    updated_date?: string
  }

  /** Карт нэмэх хүсэлт */
  interface TpayCardCreateReq {
    card_no: string
    expiry_date: string
    cvv: string
    card_holder_name: string
  }

  /** Карт баталгаажуулах хүсэлт */
  interface TpayCardConfirmReq {
    card_id: number
    otp_amount: string
  }

  /** Гүйлгээ */
  interface TpayTransaction {
    id: number
    transaction_id: string
    type: 'qr_pay' | 'p2p' | 'top_up' | 'withdraw'
    amount: number
    fee?: number
    currency: string
    status: 'pending' | 'completed' | 'failed' | 'cancelled'
    from_account?: string
    to_account?: string
    from_user_name?: string
    to_user_name?: string
    description?: string
    reference?: string
    created_date: string
    completed_date?: string
  }

  /** QR төлбөр хүсэлт */
  interface TpayQRPayReq {
    qr_string: string
    amount?: number
    pin?: string
  }

  /** P2P шилжүүлэг хүсэлт */
  interface TpayP2PReq {
    to_account: string
    amount: number
    description?: string
    pin?: string
  }

  /** Гүйлгээний хариу */
  interface TpayTransactionRes {
    transaction_id: string
    status: 'pending' | 'completed' | 'failed'
    message?: string
  }
}
