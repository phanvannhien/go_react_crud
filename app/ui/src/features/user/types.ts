export interface User {
    id: string;
    email: string;
    role: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface UserListParams {
    page?: number;
    limit?: number;
    sort?: string;
    order?: string;
    search?: string;
    role?: string;
    is_active?: boolean;
}

export interface UpdateUserRequest {
    email?: string;
    role?: string;
    is_active?: boolean;
}
