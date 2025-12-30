import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import type { BeforeResolveFunction } from '@/core/base.ts';

import { Customer, CustomerType, type CustomerInfo } from '@/models/customer.ts';

import { isEquals } from '@/lib/common.ts';
import services, { type ApiResponsePromise } from '@/lib/services.ts';
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
    function getAllCustomers(params: { visible_only?: boolean; customer_type?: CustomerType } = {}): ApiResponsePromise<CustomerInfo[]> {
        return services.get('/api/v1/customers/list.json', { params }).then(response => {
            const data = response.data;
            const customers = Customer.ofMulti(data);
            loadCustomersList(customers);
            customersListStateInvalid.value = false;
            return customers;
        }).catch(error => {
            logger.post('error', 'Failed to get customers list', error);
            throw error;
        });
    }

    function getCustomer(id: string): ApiResponsePromise<Customer> {
        return services.get('/api/v1/customers/get.json', { params: { id } }).then(response => {
            const customer = Customer.of(response.data);
            modifyCustomerInCustomersList(customer);
            return customer;
        }).catch(error => {
            logger.post('error', 'Failed to get customer', error);
            throw error;
        });
    }

    function createCustomer(data: { name: string; customer_type: CustomerType; address?: string; contacts?: string; contacts_info?: string; comment?: string; hidden?: boolean; client_session_id?: string }): ApiResponsePromise<Customer> {
        return services.post('/api/v1/customers/add.json', data).then(response => {
            const customer = Customer.of(response.data);
            addCustomerToCustomersList(customer);
            return customer;
        }).catch(error => {
            logger.post('error', 'Failed to create customer', error);
            throw error;
        });
    }

    function modifyCustomer(data: { id: string; name: string; customer_type: CustomerType; address?: string; contacts?: string; contacts_info?: string; comment?: string; hidden?: boolean }): ApiResponsePromise<Customer> {
        return services.post('/api/v1/customers/modify.json', data).then(response => {
            const customer = Customer.of(response.data);
            modifyCustomerInCustomersList(customer);
            return customer;
        }).catch(error => {
            logger.post('error', 'Failed to modify customer', error);
            throw error;
        });
    }

    function deleteCustomer(id: string): ApiResponsePromise<void> {
        return services.post('/api/v1/customers/delete.json', { id }).then(() => {
            removeCustomerFromCustomersList(id);
        }).catch(error => {
            logger.post('error', 'Failed to delete customer', error);
            throw error;
        });
    }

    function hideCustomer(id: string, hidden: boolean): ApiResponsePromise<void> {
        return services.post('/api/v1/customers/hide.json', { id, hidden }).then(() => {
            const customer = customersMap.value[id];
            if (customer) {
                customer.hidden = hidden;
            }
        }).catch(error => {
            logger.post('error', 'Failed to hide/show customer', error);
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
        getCustomer,
        createCustomer,
        modifyCustomer,
        deleteCustomer,
        hideCustomer
    };
});

export type CustomersBeforeResolveFunction = BeforeResolveFunction<ReturnType<typeof useCustomersStore>>;
