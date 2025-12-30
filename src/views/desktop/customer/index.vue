<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
                        <span>{{ tt('Customer Management') }}</span>
                        <v-btn class="ms-3" color="default" variant="outlined"
                            :disabled="loading || updating" @click="add">{{ tt('Add') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" size="24" class="ms-2" :icon="true"
                            :disabled="loading || updating" :loading="loading" @click="reload">
                            <template #loader>
                                <v-progress-circular indeterminate size="20" />
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                        <v-spacer />
                        <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                            :disabled="loading || updating" :icon="true">
                            <v-icon :icon="mdiDotsVertical" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-item :prepend-icon="mdiEyeOutline"
                                        :title="tt('Show Hidden Customers')" v-if="!showHidden"
                                        @click="showHidden = true"></v-list-item>
                                    <v-list-item :prepend-icon="mdiEyeOffOutline"
                                        :title="tt('Hide Hidden Customers')" v-if="showHidden"
                                        @click="showHidden = false"></v-list-item>
                                </v-list>
                            </v-menu>
                        </v-btn>
                    </div>
                </template>

                <v-table class="customers-table table-striped" :hover="!loading">
                    <thead>
                        <tr>
                            <th style="width: 20%;">
                                <div class="d-flex align-center">
                                    <span>{{ tt('Name') }}</span>
                                </div>
                            </th>
                            <th style="width: 15%;">
                                <div class="d-flex align-center">
                                    <span>{{ tt('Type') }}</span>
                                </div>
                            </th>
                            <th style="width: 20%;">
                                <div class="d-flex align-center">
                                    <span>{{ tt('Contacts') }}</span>
                                </div>
                            </th>
                            <th style="width: 15%;">
                                <div class="d-flex align-center">
                                    <span>{{ tt('Phone') }}</span>
                                </div>
                            </th>
                            <th style="width: 30%;">
                                <div class="d-flex align-center">
                                    <span>{{ tt('Operation') }}</span>
                                </div>
                            </th>
                        </tr>
                    </thead>

                    <tbody v-if="loading && noAvailableCustomer">
                        <tr :key="itemIdx" v-for="itemIdx in [1, 2, 3, 4, 5]">
                            <td colspan="5" class="px-0">
                                <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                            </td>
                        </tr>
                    </tbody>

                    <tbody v-if="!loading && noAvailableCustomer">
                        <tr>
                            <td colspan="5">{{ tt('No Available Customer') }}</td>
                        </tr>
                    </tbody>

                    <template v-for="customer in customers" :key="customer.id">
                        <tr class="customers-table-row text-sm" v-if="showHidden || !customer.hidden">
                            <td>
                                <div class="d-flex align-center">
                                    <v-badge class="right-bottom-icon" color="secondary" location="bottom right"
                                        offset-x="4" :icon="mdiEyeOffOutline" v-if="customer.hidden">
                                        <v-icon size="20" start :icon="mdiAccountBoxOutline" />
                                    </v-badge>
                                    <v-icon size="20" start :icon="mdiAccountBoxOutline" v-else />
                                    <span class="customer-name">{{ customer.name }}</span>
                                </div>
                            </td>
                            <td>
                                <span>{{ getCustomerTypeName(customer.customerType) }}</span>
                            </td>
                            <td>
                                <span>{{ customer.contacts || '-' }}</span>
                            </td>
                            <td>
                                <span>{{ customer.contactsInfo || '-' }}</span>
                            </td>
                            <td>
                                <div class="d-flex align-center">
                                    <v-btn class="px-2" color="default" density="comfortable" variant="text"
                                        :class="{ 'd-none': loading, 'hover-display': !loading }"
                                        :prepend-icon="customer.hidden ? mdiEyeOutline : mdiEyeOffOutline"
                                        :loading="customerHiding[customer.id]" :disabled="loading || updating"
                                        @click="hide(customer, !customer.hidden)">
                                        <template #loader>
                                            <v-progress-circular indeterminate size="20" width="2" />
                                        </template>
                                        {{ customer.hidden ? tt('Show') : tt('Hide') }}
                                    </v-btn>
                                    <v-btn class="px-2" color="default" density="comfortable" variant="text"
                                        :class="{ 'd-none': loading, 'hover-display': !loading }"
                                        :prepend-icon="mdiPencilOutline" :loading="customerUpdating[customer.id]"
                                        :disabled="loading || updating" @click="edit(customer)">
                                        <template #loader>
                                            <v-progress-circular indeterminate size="20" width="2" />
                                        </template>
                                        {{ tt('Edit') }}
                                    </v-btn>
                                    <v-btn class="px-2" color="default" density="comfortable" variant="text"
                                        :class="{ 'd-none': loading, 'hover-display': !loading }"
                                        :prepend-icon="mdiDeleteOutline" :loading="customerRemoving[customer.id]"
                                        :disabled="loading || updating" @click="remove(customer)">
                                        <template #loader>
                                            <v-progress-circular indeterminate size="20" width="2" />
                                        </template>
                                        {{ tt('Delete') }}
                                    </v-btn>
                                </div>
                            </td>
                        </tr>
                    </template>

                    <tbody v-if="newCustomer">
                        <tr class="text-sm">
                            <td colspan="5">
                                <v-text-field class="w-100" type="text" color="primary" density="compact"
                                    variant="underlined" :disabled="loading || updating"
                                    :label="tt('Name')" v-model="newCustomer.name"
                                    @keyup.enter="save(newCustomer)">
                                </v-text-field>
                            </td>
                        </tr>
                        <tr class="text-sm">
                            <td colspan="5">
                                <v-select class="customer-type-select" :items="customerTypeOptions"
                                    :label="tt('Customer Type')" v-model="newCustomer.customerType"
                                    :disabled="loading || updating" density="compact" variant="underlined"
                                    return-object />
                            </td>
                        </tr>
                        <tr class="text-sm">
                            <td colspan="5">
                                <v-text-field class="w-100" type="text" color="primary" density="compact"
                                    variant="underlined" :disabled="loading || updating"
                                    :label="tt('Address')" v-model="newCustomer.address">
                                </v-text-field>
                            </td>
                        </tr>
                        <tr class="text-sm">
                            <td colspan="5">
                                <v-text-field class="w-100" type="text" color="primary" density="compact"
                                    variant="underlined" :disabled="loading || updating"
                                    :label="tt('Contacts')" v-model="newCustomer.contacts">
                                </v-text-field>
                            </td>
                        </tr>
                        <tr class="text-sm">
                            <td colspan="5">
                                <v-text-field class="w-100" type="text" color="primary" density="compact"
                                    variant="underlined" :disabled="loading || updating"
                                    :label="tt('Phone')" v-model="newCustomer.contactsInfo">
                                </v-text-field>
                            </td>
                        </tr>
                        <tr class="text-sm">
                            <td colspan="5">
                                <v-text-field class="w-100" type="text" color="primary" density="compact"
                                    variant="underlined" :disabled="loading || updating"
                                    :label="tt('Comment')" v-model="newCustomer.comment">
                                </v-text-field>
                            </td>
                        </tr>
                        <tr class="text-sm">
                            <td colspan="5">
                                <div class="d-flex align-center">
                                    <v-btn class="px-2" density="comfortable" variant="text" :prepend-icon="mdiCheck"
                                        :loading="customerUpdating['']"
                                        :disabled="loading || updating || !isNewCustomerModified" @click="save(newCustomer)">
                                        <template #loader>
                                            <v-progress-circular indeterminate size="20" width="2" />
                                        </template>
                                        {{ tt('Save') }}
                                    </v-btn>
                                    <v-btn class="px-2" color="default" density="comfortable" variant="text"
                                        :prepend-icon="mdiClose" :disabled="loading || updating"
                                        @click="cancelSave">
                                        {{ tt('Cancel') }}
                                    </v-btn>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </v-table>
            </v-card>
        </v-col>
    </v-row>

    <customer-edit-dialog ref="editDialog" />
    <confirm-dialog ref="confirmDialog" />
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import CustomerEditDialog from './dialogs/EditDialog.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useCustomersStore } from '@/stores/customer.ts';
import { Customer, CustomerType, type CustomerInfo } from '@/models/customer.ts';
import { generateClientSessionId } from '@/lib/session.ts';

import {
    mdiRefresh,
    mdiPencilOutline,
    mdiCheck,
    mdiClose,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiDeleteOutline,
    mdiDotsVertical,
    mdiAccountBoxOutline
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type CustomerEditDialogType = InstanceType<typeof CustomerEditDialog>;

const { tt } = useI18n();

const customersStore = useCustomersStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const editDialog = useTemplateRef<CustomerEditDialogType>('editDialog');

const newCustomer = ref<Customer | null>(null);
const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const customerUpdating = ref<Record<string, boolean>>({});
const customerHiding = ref<Record<string, boolean>>({});
const customerRemoving = ref<Record<string, boolean>>({});
const showHidden = ref<boolean>(false);

const customers = computed<Customer[]>(() => customersStore.allCustomers);
const noAvailableCustomer = computed<boolean>(() => {
    if (newCustomer.value) {
        return false;
    }
    return customers.value.filter(c => showHidden.value || !c.hidden).length === 0;
});
const isNewCustomerModified = computed<boolean>(() => {
    return newCustomer.value !== null && newCustomer.value.name.trim() !== '';
});

const customerTypeOptions = computed(() => [
    { title: tt('Customer'), value: CustomerType.CUSTOMER },
    { title: tt('Supplier'), value: CustomerType.SUPPLIER },
    { title: tt('Both'), value: CustomerType.BOTH }
]);

function getCustomerTypeName(type: CustomerType): string {
    switch (type) {
        case CustomerType.CUSTOMER:
            return tt('Customer');
        case CustomerType.SUPPLIER:
            return tt('Supplier');
        case CustomerType.BOTH:
            return tt('Both');
        default:
            return '';
    }
}

function reload(): void {
    loading.value = true;

    customersStore.getAllCustomers({
        visible_only: !showHidden.value
    }).then(() => {
        loading.value = false;
        snackbar.value?.showMessage(tt('Customer list has been updated'));
    }).catch(error => {
        loading.value = false;
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function add(): void {
    newCustomer.value = Customer.createNew('', CustomerType.CUSTOMER);
}

function edit(customer: Customer): void {
    editDialog.value?.open(customer).then((modified) => {
        if (modified) {
            reload();
        }
    });
}

function save(customer: Customer): void {
    if (!customer) {
        return;
    }

    updating.value = true;
    customerUpdating.value[customer.id || ''] = true;

    const clientSessionId = generateClientSessionId();

    customersStore.createCustomer({
        name: customer.name,
        customer_type: customer.customerType,
        address: customer.address,
        contacts: customer.contacts,
        contacts_info: customer.contactsInfo,
        comment: customer.comment,
        hidden: customer.hidden,
        client_session_id: clientSessionId
    }).then(() => {
        updating.value = false;
        customerUpdating.value[customer.id || ''] = false;
        newCustomer.value = null;
    }).catch(error => {
        updating.value = false;
        customerUpdating.value[customer.id || ''] = false;
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function cancelSave(): void {
    newCustomer.value = null;
}

function hide(customer: Customer, hidden: boolean): void {
    updating.value = true;
    customerHiding.value[customer.id] = true;

    customersStore.hideCustomer(customer.id, hidden).then(() => {
        updating.value = false;
        customerHiding.value[customer.id] = false;
    }).catch(error => {
        updating.value = false;
        customerHiding.value[customer.id] = false;
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function remove(customer: Customer): void {
    confirmDialog.value?.open(tt('Are you sure you want to delete this customer?')).then(() => {
        updating.value = true;
        customerRemoving.value[customer.id] = true;

        customersStore.deleteCustomer(customer.id).then(() => {
            updating.value = false;
            customerRemoving.value[customer.id] = false;
        }).catch(error => {
            updating.value = false;
            customerRemoving.value[customer.id] = false;
            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

watch(showHidden, () => {
    reload();
});

// Load initial data
customersStore.getAllCustomers({
    visible_only: false
}).then(() => {
    loading.value = false;
}).catch(error => {
    loading.value = false;
    if (!error.processed) {
        snackbar.value?.showError(error);
    }
});
</script>

<style scoped>
.customers-table tr.customers-table-row .hover-display {
    display: none;
}

.customers-table tr.customers-table-row:hover .hover-display {
    display: inline-grid;
}

.customers-table tr:not(:last-child)>td>div {
    padding-bottom: 1px;
}

.customers-table tr.customers-table-row .right-bottom-icon .v-badge__badge {
    padding-bottom: 1px;
}

.customers-table .customer-name {
    font-size: 0.875rem;
}

.customers-table .customer-type-select {
    font-size: 0.875rem;
}
</style>
