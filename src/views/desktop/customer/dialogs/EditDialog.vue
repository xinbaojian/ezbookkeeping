<template>
    <v-dialog :model-value="isOpen" :max-width="600" @update:model-value="onDialogStateChange" persistent>
        <v-card>
            <v-card-title>{{ tt('Edit Customer') }}</v-card-title>
            <v-card-text>
                <v-text-field
                    v-model="editingCustomer.name"
                    :label="tt('Name')"
                    :disabled="updating"
                    density="compact"
                    variant="underlined"
                    class="mb-4"
                />

                <v-select
                    v-model="editingCustomer.customerType"
                    :items="customerTypeOptions"
                    :label="tt('Customer Type')"
                    :disabled="updating"
                    density="compact"
                    variant="underlined"
                    return-object
                    class="mb-4"
                />

                <v-text-field
                    v-model="editingCustomer.address"
                    :label="tt('Address')"
                    :disabled="updating"
                    density="compact"
                    variant="underlined"
                    class="mb-4"
                />

                <v-text-field
                    v-model="editingCustomer.contacts"
                    :label="tt('Contacts')"
                    :disabled="updating"
                    density="compact"
                    variant="underlined"
                    class="mb-4"
                />

                <v-text-field
                    v-model="editingCustomer.contactsInfo"
                    :label="tt('Phone')"
                    :disabled="updating"
                    density="compact"
                    variant="underlined"
                    class="mb-4"
                />

                <v-text-field
                    v-model="editingCustomer.comment"
                    :label="tt('Comment')"
                    :disabled="updating"
                    density="compact"
                    variant="underlined"
                    class="mb-4"
                />

                <v-checkbox
                    v-model="editingCustomer.hidden"
                    :label="tt('Hidden')"
                    :disabled="updating"
                    density="compact"
                />
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn
                    color="default"
                    variant="text"
                    :disabled="updating"
                    @click="cancel"
                >
                    {{ tt('Cancel') }}
                </v-btn>
                <v-btn
                    color="primary"
                    variant="tonal"
                    :loading="updating"
                    :disabled="!isModified"
                    @click="save"
                >
                    {{ tt('Save') }}
                </v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useI18n } from '@/locales/helpers.ts';
import { useCustomersStore } from '@/stores/customer.ts';
import { Customer, CustomerType } from '@/models/customer.ts';

const { tt } = useI18n();
const customersStore = useCustomersStore();

const isOpen = ref<boolean>(false);
const updating = ref<boolean>(false);
const originalCustomer = ref<Customer | null>(null);
const editingCustomer = ref<Customer>(Customer.createNew('', CustomerType.CUSTOMER));

const customerTypeOptions = computed(() => [
    { title: tt('Customer'), value: CustomerType.CUSTOMER },
    { title: tt('Supplier'), value: CustomerType.SUPPLIER },
    { title: tt('Both'), value: CustomerType.BOTH }
]);

const isModified = computed<boolean>(() => {
    if (!originalCustomer.value) {
        return false;
    }
    return !editingCustomer.value.equals(originalCustomer.value);
});

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
    (e: 'save', modified: boolean): void;
});

function onDialogStateChange(value: boolean): void {
    if (!value && !updating.value) {
        isOpen.value = false;
        emit('update:modelValue', false);
    }
}

function open(customer: Customer): Promise<boolean> {
    return new Promise<boolean>((resolve) => {
        originalCustomer.value = customer;
        editingCustomer.value = Customer.of({
            id: customer.id,
            name: customer.name,
            customer_type: customer.customerType,
            address: customer.address,
            contacts: customer.contacts,
            contacts_info: customer.contactsInfo,
            comment: customer.comment,
            hidden: customer.hidden,
            created_time: customer.createdTime.toISOString(),
            updated_time: customer.updatedTime.toISOString()
        });
        isOpen.value = true;
        emit('update:modelValue', true);

        // Set up a watcher to resolve when dialog closes
        const stopWatch = watch(isOpen, (newValue) => {
            if (!newValue) {
                stopWatch();
                resolve(false);
            }
        });
    });
}

function cancel(): void {
    isOpen.value = false;
    emit('update:modelValue', false);
}

function save(): void {
    if (!originalCustomer.value || !isModified.value) {
        return;
    }

    updating.value = true;

    customersStore.modifyCustomer({
        id: editingCustomer.value.id,
        name: editingCustomer.value.name,
        customer_type: editingCustomer.value.customerType,
        address: editingCustomer.value.address,
        contacts: editingCustomer.value.contacts,
        contacts_info: editingCustomer.value.contactsInfo,
        comment: editingCustomer.value.comment,
        hidden: editingCustomer.value.hidden
    }).then(() => {
        updating.value = false;
        isOpen.value = false;
        emit('update:modelValue', false);
        emit('save', true);
    }).catch(error => {
        updating.value = false;
        if (!error.processed) {
            // Show error - caller should handle this
            console.error(error);
        }
    });
}

defineExpose({
    open
});
</script>
