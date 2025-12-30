import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { Customer, CustomerType } from '@/models/customer.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useCustomersStore = defineStore('customers', () => {
    const allCustomers = ref<Customer[]>([]);
    const customersMap = ref<Record<string, Customer>>({});
    const customersListStateInvalid = ref<boolean>(true);

    const allAvailableCustomersCount = computed<number>(() => {
        return allCustomers.value.filter(c => !c.hidden).length;
    });

    const hasVisibleCustomers = computed<boolean>(() => {
        return allCustomers.value.some(c => !c.hidden);
    });

    function loadCustomersList(customers: Customer[]): void {
        allCustomers.value = customers;
        customersMap.value = {};

        for (const customer of customers) {
            customersMap.value[customer.id] = customer;
        }
    }

    function addCustomerToCustomersList(customer: Customer): void {
        allCustomers.value.unshift(customer);
        customersMap.value[customer.id] = customer;
    }

    function modifyCustomerInCustomersList(customer: Customer): void {
        const existingCustomer = customersMap.value[customer.id];

        if (existingCustomer) {
            const index = allCustomers.value.findIndex(c => c.id === customer.id);

            if (index >= 0) {
                allCustomers.value[index] = customer;
            }

            customersMap.value[customer.id] = customer;
        }
    }

    function removeCustomerFromCustomersList(customerId: string): void {
        const index = allCustomers.value.findIndex(c => c.id === customerId);

        if (index >= 0) {
            allCustomers.value.splice(index, 1);
        }

        delete customersMap.value[customerId];
    }

    function getCustomerById(customerId: string): Customer | undefined {
        return customersMap.value[customerId];
    }

    function getCustomersByType(customerType: CustomerType): Customer[] {
        return allCustomers.value.filter(c => c.customerType === customerType);
    }

    // API calls
    function getAllCustomers(params: { visible_only?: boolean; customer_type?: CustomerType } = {}): Promise<Customer[]> {
        return services.getAllCustomers(params).then(response => {
            const data = response.data;
            if (!data || !data.success || !data.result) {
                throw new Error('Unable to retrieve customers list');
            }
            const customers = Customer.ofMulti(data.result);
            loadCustomersList(customers);
            customersListStateInvalid.value = false;
            return customers;
        }).catch(error => {
            logger.error('Failed to get customers list', error);
            throw error;
        });
    }

    function getAllCustomersWithPagination(params: { visible_only?: boolean; customer_type?: CustomerType; page?: number; page_size?: number } = {}): Promise<{ customers: Customer[]; total: number; totalPages: number }> {
        return services.getAllCustomersWithPagination(params).then(response => {
            const data = response.data;
            if (!data || !data.success || !data.result) {
                throw new Error('Unable to retrieve customers list');
            }
            const customers = Customer.ofMulti(data.result.customers);
            loadCustomersList(customers);
            customersListStateInvalid.value = false;
            return {
                customers,
                total: data.result.total,
                totalPages: data.result.total_pages
            };
        }).catch(error => {
            logger.error('Failed to get customers list', error);
            throw error;
        });
    }

    function getCustomer(id: string): Promise<Customer> {
        return services.getCustomer({ id }).then(response => {
            const data = response.data;
            if (!data || !data.success || !data.result) {
                throw new Error('Unable to retrieve customer');
            }
            const customer = Customer.of(data.result);
            modifyCustomerInCustomersList(customer);
            return customer;
        }).catch(error => {
            logger.error('Failed to get customer', error);
            throw error;
        });
    }

    function createCustomer(data: { name: string; customer_type: CustomerType; address?: string; contacts?: string; contacts_info?: string; comment?: string; hidden?: boolean; client_session_id?: string }): Promise<Customer> {
        return services.createCustomer(data).then(response => {
            const responseData = response.data;
            if (!responseData || !responseData.success || !responseData.result) {
                throw new Error('Unable to create customer');
            }
            const customer = Customer.of(responseData.result);
            addCustomerToCustomersList(customer);
            return customer;
        }).catch(error => {
            logger.error('Failed to create customer', error);
            throw error;
        });
    }

    function modifyCustomer(data: { id: string; name: string; customer_type: CustomerType; address?: string; contacts?: string; contacts_info?: string; comment?: string; hidden?: boolean }): Promise<Customer> {
        return services.modifyCustomer(data).then(response => {
            const responseData = response.data;
            if (!responseData || !responseData.success || !responseData.result) {
                throw new Error('Unable to modify customer');
            }
            const customer = Customer.of(responseData.result);
            modifyCustomerInCustomersList(customer);
            return customer;
        }).catch(error => {
            logger.error('Failed to modify customer', error);
            throw error;
        });
    }

    function deleteCustomer(id: string): Promise<void> {
        return services.deleteCustomer({ id }).then(() => {
            removeCustomerFromCustomersList(id);
        }).catch(error => {
            logger.error('Failed to delete customer', error);
            throw error;
        });
    }

    function hideCustomer(id: string, hidden: boolean): Promise<void> {
        return services.hideCustomer({ id, hidden }).then(() => {
            const customer = customersMap.value[id];
            if (customer) {
                customer.hidden = hidden;
            }
        }).catch(error => {
            logger.error('Failed to hide/show customer', error);
            throw error;
        });
    }

    return {
        allCustomers,
        customersMap,
        customersListStateInvalid,
        allAvailableCustomersCount,
        hasVisibleCustomers,
        loadCustomersList,
        addCustomerToCustomersList,
        modifyCustomerInCustomersList,
        removeCustomerFromCustomersList,
        getCustomerById,
        getCustomersByType,
        getAllCustomers,
        getAllCustomersWithPagination,
        getCustomer,
        createCustomer,
        modifyCustomer,
        deleteCustomer,
        hideCustomer
    };
});
