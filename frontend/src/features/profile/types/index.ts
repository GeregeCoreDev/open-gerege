export interface UserProfile {
    id: string;
    first_name: string;
    last_name: string;
    email: string;
    phone_no: string;
    profile_img_url?: string;
    is_org: boolean;
    reg_no?: string;
    // Add other fields observed in Swagger if needed
}

export interface UserState {
    user: UserProfile | null;
    isLoading: boolean;
    error: string | null;
    fetchProfile: () => Promise<void>;
    logout: () => void;
}
