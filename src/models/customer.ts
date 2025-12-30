// Customer types
export enum CustomerType {
    CUSTOMER = 1,
    SUPPLIER = 2,
    BOTH = 3
}

export interface CustomerInfo {
    id: string;
    name: string;
    customer_type: CustomerType;
    address: string;
    contacts: string;
    contacts_info: string;
    comment: string;
    hidden: boolean;
    created_time: string;
    updated_time: string;
}

export interface CustomerListResponse {
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
    customers: CustomerInfo[];
}

export interface CustomerCreateRequest {
    client_session_id?: string;
    name: string;
    customer_type: CustomerType;
    address?: string;
    contacts?: string;
    contacts_info?: string;
    comment?: string;
    hidden?: boolean;
}

export interface CustomerModifyRequest {
    id: string;
    name: string;
    customer_type: CustomerType;
    address?: string;
    contacts?: string;
    contacts_info?: string;
    comment?: string;
    hidden?: boolean;
}

export interface CustomerDeleteRequest {
    id: string;
}

export interface CustomerHideRequest {
    id: string;
    hidden: boolean;
}

export interface CustomerListRequest {
    visible_only?: boolean;
    customer_type?: CustomerType;
    page?: number;
    page_size?: number;
}

export class Customer {
    public id: string;
    public name: string;
    public customerType: CustomerType;
    public address: string;
    public contacts: string;
    public contactsInfo: string;
    public comment: string;
    public hidden: boolean;
    public createdTime: Date;
    public updatedTime: Date;

    private constructor(info: CustomerInfo) {
        this.id = info.id;
        this.name = info.name;
        this.customerType = info.customer_type;
        this.address = info.address;
        this.contacts = info.contacts;
        this.contactsInfo = info.contacts_info;
        this.comment = info.comment;
        this.hidden = info.hidden;
        this.createdTime = new Date(info.created_time);
        this.updatedTime = new Date(info.updated_time);
    }

    public get visible(): boolean {
        return !this.hidden;
    }

    public get customerTypeName(): string {
        switch (this.customerType) {
            case CustomerType.CUSTOMER:
                return 'Customer';
            case CustomerType.SUPPLIER:
                return 'Supplier';
            case CustomerType.BOTH:
                return 'Both';
            default:
                return '';
        }
    }

    public equals(other: Customer): boolean {
        return this.id === other.id &&
            this.name === other.name &&
            this.customerType === other.customerType &&
            this.address === other.address &&
            this.contacts === other.contacts &&
            this.contactsInfo === other.contactsInfo &&
            this.comment === other.comment &&
            this.hidden === other.hidden;
    }

    public toModifyRequest(): CustomerModifyRequest {
        return {
            id: this.id,
            name: this.name,
            customer_type: this.customerType,
            address: this.address,
            contacts: this.contacts,
            contacts_info: this.contactsInfo,
            comment: this.comment,
            hidden: this.hidden
        };
    }

    public static of(info: CustomerInfo): Customer {
        return new Customer(info);
    }

    public static ofMulti(infos: CustomerInfo[]): Customer[] {
        return infos.map(info => Customer.of(info));
    }

    public static createNew(name: string, customerType: CustomerType): Customer {
        return new Customer({
            id: '',
            name: name,
            customer_type: customerType,
            address: '',
            contacts: '',
            contacts_info: '',
            comment: '',
            hidden: false,
            created_time: new Date().toISOString(),
            updated_time: new Date().toISOString()
        });
    }
}
