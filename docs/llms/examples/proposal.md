# {{ proposal_title }}

**Prepared for:** {{ client.name }}
**Prepared by:** {{ company.name }}
**Date:** {{ proposal_date }}
**Valid Until:** {{ valid_until }}

---

{{@toc}}

---

## Introduction

Dear {{ client.contact_name }},

{{ introduction }}

We are excited to present this proposal outlining how {{ company.name }} can help {{ client.name }} achieve its goals.

{{@pagebreak}}

## Understanding Your Needs

{{ needs_assessment }}

### Key Challenges

{{#challenges}}
- {{ . }}
{{/challenges}}

### Success Criteria

{{#success_criteria}}
- {{ . }}
{{/success_criteria}}

---

## Proposed Solution

{{@highlight}}
{{ solution_overview }}
{{/highlight}}

{{@columns}}
**Our Approach**

{{ approach_description }}

We bring {{ company.years_experience }} years of experience and a proven track record in delivering similar solutions.

---

**Key Differentiators**

{{#differentiators}}
- {{ . }}
{{/differentiators}}
{{/columns}}

{{@pagebreak}}

## Scope of Work

{{#phases}}
### Phase {{ number }}: {{ name }}

{{ description }}

**Deliverables:**
{{#deliverables}}
- {{ . }}
{{/deliverables}}

**Duration:** {{ duration }}

{{/phases}}

---

## Timeline

| Phase | Start | End | Milestone |
|-------|-------|-----|-----------|
{{#timeline}}
| {{ phase }} | {{ start }} | {{ end }} | {{ milestone }} |
{{/timeline}}

{{@note}}
Timeline assumes project kickoff within 2 weeks of proposal acceptance. Dates will be adjusted accordingly upon signed agreement.
{{/note}}

{{@pagebreak}}

## Investment

{{@columns}}
### Project Fees

| Component | Investment |
|-----------|-----------|
{{#pricing}}
| {{ item }} | ${{ amount }} |
{{/pricing}}

---

### Payment Schedule

{{#payment_schedule}}
- **{{ milestone }}:** {{ percentage }}% (${{ amount }})
{{/payment_schedule}}
{{/columns}}

---

**Total Investment: ${{ total_investment }}**

{{#discount}}
{{@highlight}}
**Early Signing Bonus:** Sign by {{ discount.deadline }} and receive {{ discount.percentage }}% off (${{ discount.savings }} savings).
{{/highlight}}
{{/discount}}

---

## Why {{ company.name }}

{{ company_overview }}

### Relevant Experience

{{#case_studies}}
**{{ client }}** - {{ industry }}

{{ description }}

*Result: {{ result }}*

{{/case_studies}}

---

## Next Steps

{{#next_steps}}
1. {{ . }}
{{/next_steps}}

{{@pagebreak}}

## Terms and Conditions

{{ terms_summary }}

{{@warning}}
This proposal is valid until {{ valid_until }}. Pricing and availability are subject to change after this date.
{{/warning}}

---

## Acceptance

To proceed with this proposal, please sign below and return to {{ company.contact_email }}.

**Client Signature:** _________________________ **Date:** _____________

**Printed Name:** {{ client.contact_name }}

**Title:** _________________________ **Company:** {{ client.name }}

---

*We look forward to partnering with {{ client.name }} on this exciting initiative.*

**{{ company.name }}**
{{ company.address }}
{{ company.phone }} | {{ company.contact_email }}
