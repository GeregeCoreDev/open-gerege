import { apiClient } from '@/lib/api-client';
import { UserProfile } from './types';

export const profileApi = {
    getProfile: () => apiClient.get<UserProfile>('/me/profile'),
};
