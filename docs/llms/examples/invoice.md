# Invoice

**Invoice #{{ invoice_number }}**

**Date:** {{ invoice_date }}
**Due Date:** {{ due_date }}

---

## From

**{{ company.name }}**
{{ company.address.street }}
{{ company.address.city }}, {{ company.address.state }} {{ company.address.zip }}
{{ company.email }}

---

## Bill To

**{{ client.name }}**
{{ client.address.street }}
{{ client.address.city }}, {{ client.address.state }} {{ client.address.zip }}
{{#client.email}}
{{ client.email }}
{{/client.email}}

---

## Line Items

| Description | Quantity | Unit Price | Amount |
|-------------|----------|------------|--------|
{{#line_items}}
| {{ description }} | {{ quantity }} | ${{ unit_price }} | ${{ amount }} |
{{/line_items}}

---

{{@columns}}
**Subtotal:** ${{ subtotal }}

{{#discount}}
**Discount ({{ discount.description }}):** -${{ discount.amount }}
{{/discount}}

**Tax ({{ tax_rate }}%):** ${{ tax_amount }}

---

**Total Due:** ${{ total }}
{{/columns}}

---

{{#notes}}
{{@note}}
{{ notes }}
{{/note}}
{{/notes}}

{{#payment_instructions}}
## Payment Instructions

{{ payment_instructions }}
{{/payment_instructions}}

{{^paid}}
{{@warning}}
This invoice is unpaid. Please remit payment by {{ due_date }}.
{{/warning}}
{{/paid}}

{{#paid}}
{{@highlight}}
**Payment Received** - Thank you for your payment on {{ payment_date }}.
{{/highlight}}
{{/paid}}

---

*Thank you for your business!*
