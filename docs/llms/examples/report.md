# {{ report_title }}

**Prepared by:** {{ author.name }}, {{ author.title }}
**Date:** {{ report_date }}
**Client:** {{ client.name }}

---

{{@toc}}

---

## Executive Summary

{{@highlight}}
{{ executive_summary }}
{{/highlight}}

{{@pagebreak}}

## Introduction

{{ introduction }}

### Project Background

{{ project_background }}

### Objectives

{{#objectives}}
- {{ . }}
{{/objectives}}

---

## Methodology

{{ methodology.overview }}

### Data Collection

{{#methodology.data_sources}}
- **{{ name }}:** {{ description }}
{{/methodology.data_sources}}

### Analysis Approach

{{ methodology.analysis }}

{{@pagebreak}}

## Findings

{{#findings}}
### {{ title }}

{{ description }}

{{#data}}
| Metric | Value | Change |
|--------|-------|--------|
{{#metrics}}
| {{ name }} | {{ value }} | {{ change }} |
{{/metrics}}
{{/data}}

{{#is_positive}}
{{@note}}
This finding indicates positive progress toward stated objectives.
{{/note}}
{{/is_positive}}

{{#is_critical}}
{{@warning}}
This finding requires immediate attention and follow-up action.
{{/warning}}
{{/is_critical}}

{{/findings}}

{{@pagebreak}}

## Recommendations

Based on our analysis, we recommend the following actions:

{{#recommendations}}
### {{ priority }}. {{ title }}

{{ description }}

**Timeline:** {{ timeline }}
**Estimated Impact:** {{ impact }}

{{/recommendations}}

---

## Conclusion

{{ conclusion }}

{{@highlight}}
**Next Steps:** {{ next_steps }}
{{/highlight}}

---

## Appendix

{{#appendix}}
### {{ title }}

{{ content }}

{{/appendix}}

---

*Report prepared by {{ author.name }} | {{ author.email }} | {{ report_date }}*
