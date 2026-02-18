# STAGE 4: SPECIALIZED SYSTEMS -- INFRASTRUCTURE ADVANCED

## Advanced Reference Documents for Power, Network, and Environmental Infrastructure

**Document ID:** STAGE4-INFRA-ADVANCED
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Specialized Systems -- These articles provide deep technical reference for the physical infrastructure that keeps the institution alive. They assume mastery of Stage 2 philosophy and Stage 3 operational procedures. They are designed for the operator who must build, maintain, and eventually replace every physical system the institution depends on.

---

## How to Read This Document

This document contains five specialized articles that belong to Stage 4 of the holm.chat Documentation Institution. Stage 2 established philosophy. Stage 3 established procedures. Stage 4 provides the deep technical knowledge required to design, build, operate, and maintain specific systems over a multi-decade horizon.

These are engineering references. They contain calculations, specifications, degradation curves, replacement schedules, and failure analysis. They are written for the operator who must make concrete decisions about concrete systems: how many solar panels, what size battery bank, which network switch, what temperature range, how many watts.

If you are building the institution for the first time, read these articles before purchasing hardware. The decisions made during initial construction determine maintenance requirements for years to come. A poorly sized solar array or an inadequately ventilated equipment room will generate compound problems that grow more expensive to fix with each passing year.

If you are maintaining an established institution, these articles serve as reference during scheduled maintenance, replacement planning, and troubleshooting. The degradation curves and replacement schedules are designed to be consulted repeatedly, not read once and forgotten.

If something in these articles no longer reflects current technology -- because solar panel efficiency has improved, or battery chemistry has changed, or network hardware has evolved -- adapt the principles to your reality. The physics of power generation, energy storage, heat dissipation, and signal integrity do not change. The specific products and configurations will.

---

---

# D4-002 -- Solar Power System Design and Maintenance

**Document ID:** D4-002
**Domain:** 4 -- Infrastructure & Power
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, OPS-001, D4-001
**Depended Upon By:** D4-003 (Battery Systems), D4-006 (Power Distribution), D4-005 (Environmental Controls). All articles involving power generation or energy independence.

---

## 1. Purpose

This article provides a complete engineering reference for the solar power system that supplies electricity to the holm.chat Documentation Institution. It covers every aspect of the system from initial design through twenty years of operation: panel selection and sizing, mounting and orientation, charge controller configuration, inverter selection, daily monitoring procedures, seasonal adjustment protocols, degradation tracking, and the full maintenance schedule.

The institution is off-grid by mandate (CON-001, Section 3.3). Off-grid means there is no utility connection to fall back on. The solar power system is not a supplement to grid power. It is the primary power source. When it fails, the institution runs on stored energy alone, and when stored energy is depleted, the institution goes dark. This reality governs every design decision in this article. There is no margin for optimistic assumptions. Every calculation includes safety factors. Every component has a maintenance schedule. Every failure mode has a documented response.

This article is written for the operator who must size a solar array from first principles, install it correctly, and keep it producing power for decades. It does not assume prior expertise in solar engineering. It does assume the willingness to learn, to measure, and to maintain.

## 2. Scope

This article covers:

- Load analysis: determining how much power the institution requires.
- Solar resource assessment: determining how much sunlight is available at the installation site.
- Panel selection: choosing panel technology, wattage, and quantity.
- Array configuration: series vs. parallel wiring, string sizing, voltage considerations.
- Mounting systems: fixed vs. adjustable tilt, ground mount vs. roof mount, wind load considerations.
- Charge controller selection and configuration: MPPT vs. PWM, sizing, programming parameters.
- Inverter selection: pure sine wave requirements, sizing for peak and continuous loads, efficiency considerations.
- System integration: how the solar array connects to the battery bank (D4-003) and the power distribution system (D4-006).
- Daily and weekly monitoring procedures.
- Seasonal adjustment protocols.
- Panel degradation tracking over twenty years.
- The complete twenty-year maintenance schedule.
- Troubleshooting common failures.

This article does not cover battery systems in depth (see D4-003), internal power distribution (see D4-006), or generator backup systems (see D4-006). It covers the solar generation side of the power system, from sunlight striking the panel to DC power delivered to the charge controller.

## 3. Background

### 3.1 Why Solar

Solar is the natural choice for an off-grid institution in most temperate and tropical climates for three reasons. First, it has no fuel cost. Once panels are installed, electricity is generated from sunlight at zero marginal cost for the remaining life of the panels. Second, solar panels have no moving parts. A well-installed panel has a twenty-five to thirty-year productive lifespan with minimal maintenance. Third, solar technology is mature, well-understood, and widely available. Replacement panels can be sourced decades after the original installation.

The primary limitation of solar is its intermittency. Panels produce power only during daylight hours, and production varies with weather, season, and latitude. This limitation is addressed through battery storage (D4-003) and, where necessary, generator backup (D4-006). The solar system must be designed with the worst-case season in mind, not the best-case season.

### 3.2 The Design Philosophy

The institution's solar system is designed according to three principles that derive from the root documents:

**Principle of sufficiency:** The system must produce enough energy to power the institution through the worst solar month of the year, accounting for panel degradation, battery losses, and a reasonable safety margin. "Enough" is defined by the load analysis, not by what seems like a comfortable number.

**Principle of maintainability:** Every component of the system must be accessible for inspection, cleanable, replaceable by the operator without specialized equipment, and documented thoroughly. A system the operator cannot maintain is a system with a countdown timer.

**Principle of graceful degradation:** When components fail -- and over twenty years, components will fail -- the system must degrade gracefully rather than catastrophically. This means redundancy in panels, modular design in the array, and clear procedures for operating at reduced capacity while repairs are made.

## 4. System Model

### 4.1 Load Analysis

Every solar system design begins with the load analysis. This determines how much energy the institution consumes, which determines how much energy the solar array must produce.

The load analysis is performed in four steps:

**Step 1: Inventory all electrical loads.** List every device in the institution that consumes electricity. For each device, record its rated wattage (from the nameplate or documentation), its typical operating wattage (measured with a watt meter -- nameplate ratings are often higher than actual consumption), and its daily usage in hours.

**Step 2: Calculate daily energy consumption.** For each device, multiply its operating wattage by its daily hours of use. Sum all devices. The result is the institution's daily energy consumption in watt-hours (Wh). A typical computing-focused institution with servers, storage arrays, networking equipment, lighting, and environmental controls consumes between 2,000 and 8,000 Wh per day, though this varies enormously with the specific equipment deployed.

**Step 3: Apply the system loss factor.** Not all energy generated by the panels reaches the loads. Energy is lost in the charge controller (typically 2-5% for MPPT), in battery charging and discharging (typically 10-20% depending on chemistry; see D4-003), in the inverter (typically 5-15%), and in wiring (typically 2-5%). The cumulative system loss factor is typically 20-35%. Divide the daily energy consumption by (1 minus the loss factor) to determine the daily energy that must be produced by the panels. If daily consumption is 5,000 Wh and the loss factor is 0.25, the panels must produce 5,000 / 0.75 = 6,667 Wh per day.

**Step 4: Apply the autonomy safety factor.** The system must be sized not for the average day but for the worst reasonable day during the worst solar month. Additionally, panels degrade over time (see Section 4.7). Apply a safety factor of 1.25 to 1.5 to the Step 3 result to account for degradation over the system's life, unusual weather events, and future load growth. Using the example above: 6,667 Wh x 1.3 = 8,667 Wh daily production required from the array.

### 4.2 Solar Resource Assessment

The solar resource at the installation site determines how many panels are needed to produce the required daily energy. The key metric is "peak sun hours" (PSH) -- the number of hours per day during which solar irradiance averages 1,000 watts per square meter. PSH varies by latitude, season, weather patterns, and local shading.

For the worst solar month (typically December or January in the Northern Hemisphere, June or July in the Southern Hemisphere), determine the average daily PSH. This data is available from solar irradiance databases, historical weather records, or direct measurement with a pyranometer over at least one full year. For temperate Northern Hemisphere locations, worst-month PSH typically ranges from 1.5 to 3.5 hours. For Southern Hemisphere or tropical locations, worst-month PSH may range from 3.0 to 5.5 hours.

To calculate required panel capacity: divide the daily production requirement (from Step 4 above) by the worst-month PSH. Using the example: 8,667 Wh / 2.5 PSH = 3,467 watts of panel capacity required.

### 4.3 Panel Selection

Panels are selected based on five criteria:

**Efficiency:** Higher-efficiency panels produce more watts per square meter, reducing the required mounting area. As of the institution's founding, monocrystalline panels offer the best efficiency (20-23%), followed by polycrystalline (15-18%). Choose monocrystalline for space-constrained installations.

**Durability:** Panels must survive twenty-five years of exposure to UV radiation, thermal cycling, wind, rain, hail, and snow. Look for IEC 61215 certification, a minimum 25-year performance warranty, and documented hail resistance ratings.

**Temperature coefficient:** Panels lose efficiency as temperature rises. The temperature coefficient, expressed as a percentage loss per degree Celsius above 25 degrees C, matters more in hot climates. Typical values range from -0.3% to -0.5% per degree C. Lower (closer to zero) is better.

**Physical dimensions and weight:** Panels must be physically mountable on the available structure. Consider the mounting system's load capacity, wind load in the local environment, and the operator's ability to handle the panels during installation and replacement.

**Availability and replaceability:** Choose panels from manufacturers with a track record of long-term availability. The institution must be able to source replacement panels for decades. Standard sizes (approximately 1.0m x 1.7m for residential panels) are more likely to remain available than unusual formats.

### 4.4 Array Configuration

Panels are wired in series (increasing voltage), parallel (increasing current), or combinations of both. The configuration must match the charge controller's input requirements.

**Series strings:** Panels wired in series add their voltages. A string of 10 panels rated at 40V open-circuit produces 400V open-circuit. Series strings must not exceed the charge controller's maximum input voltage. All panels in a series string must be the same model and orientation, because the weakest panel limits the entire string's current.

**Parallel connections:** Strings wired in parallel add their currents. Two strings each producing 10A produce 20A in parallel. The total current must not exceed the charge controller's maximum input current.

**String sizing:** Calculate the string size by dividing the charge controller's maximum input voltage by the panel's open-circuit voltage at the coldest expected temperature (voltage increases as temperature decreases -- use the temperature coefficient to calculate the cold-temperature voltage). Leave a 10% margin below the controller's maximum.

### 4.5 Charge Controller Selection

The charge controller regulates the voltage and current from the solar array to properly charge the battery bank. MPPT (Maximum Power Point Tracking) controllers are required for this institution. PWM (Pulse Width Modulation) controllers are simpler and cheaper but waste 15-30% of available power. Over twenty years, that waste is unacceptable.

Size the charge controller for the array's maximum power output plus a 20% margin. Ensure the controller supports the battery chemistry in use (see D4-003). Program the charge controller with the battery manufacturer's recommended charge parameters: bulk voltage, absorption voltage, float voltage, absorption time, and temperature compensation. Document these settings. Incorrect charge parameters are the leading cause of premature battery failure.

### 4.6 Inverter Selection

The inverter converts DC power from the battery bank to AC power for the institution's loads. Requirements:

**Pure sine wave:** Computing equipment requires clean power. Modified sine wave inverters produce harmonics that can damage sensitive electronics, cause data corruption, and generate audible noise in power supplies. Pure sine wave only.

**Continuous power rating:** Must exceed the institution's peak continuous load by at least 25%. Measure the actual peak load, do not estimate it.

**Surge capacity:** Must handle startup surges from motors and capacitive loads. Typical surge requirement is 2x the continuous rating for at least 5 seconds.

**Efficiency:** Inverter efficiency determines how much battery energy is lost in conversion. Target 93% or better at typical operating loads. Check efficiency at partial load as well -- many inverters are less efficient at light loads.

**Idle consumption:** When the institution's loads are minimal (e.g., overnight with only storage drives spinning), the inverter still consumes power. Lower idle consumption preserves battery capacity.

### 4.7 Degradation Tracking

Solar panels degrade over time. Typical degradation is 0.5-0.7% per year for monocrystalline panels, with a slightly faster initial degradation in the first year (1-3%, called Light-Induced Degradation, or LID). After twenty years, a panel typically produces 85-90% of its original rated output.

Track degradation by recording the array's daily peak output on clear days each month. Compare to the theoretical maximum (panel rating x number of panels x estimated irradiance at time of measurement). Plot this data over time. The resulting curve should show a slow, steady decline. Sudden drops indicate a component failure, not normal degradation.

If degradation exceeds 1% per year consistently, investigate: soiling, shading from vegetation growth, connection corrosion, bypass diode failure, or cell-level defects.

## 5. Rules & Constraints

- **R-D4-002-01:** The solar array must be sized for the worst solar month, not the average or best month. Optimistic sizing is system-level recklessness.
- **R-D4-002-02:** All charge controller parameters must be documented in the operational log and verified quarterly. Incorrect charge parameters void this article's degradation projections.
- **R-D4-002-03:** Panel output must be recorded at least monthly, using consistent measurement conditions, to track degradation trends. Records must be maintained for the life of the installation.
- **R-D4-002-04:** Only pure sine wave inverters are permitted for loads that include computing equipment. This constraint is absolute and non-negotiable.
- **R-D4-002-05:** The array must be designed for modular expansion. If load growth exceeds original projections, it must be possible to add panels without redesigning the entire system.
- **R-D4-002-06:** All wiring must use UV-resistant, outdoor-rated cable with appropriately sized conductors for the current carried. Undersized wiring is a fire hazard and an energy waste.
- **R-D4-002-07:** A physical disconnect switch must be installed between the array and the charge controller, accessible to the operator, to allow safe maintenance of either component.

## 6. Failure Modes

- **Panel failure.** Individual panels can fail due to cell cracking, bypass diode failure, junction box corrosion, or delamination. Impact: reduced array output proportional to the failed panel's contribution. Detection: monthly output monitoring shows unexpected drop. A single panel failure in a parallel-configured array reduces output incrementally. A panel failure in a series string can reduce the entire string's output if the bypass diode also fails.

- **Charge controller failure.** The charge controller is a single point of failure in most systems. If it fails, no solar energy reaches the batteries. Impact: total loss of charging. The institution runs on battery reserves alone. Detection: batteries stop charging despite adequate sunlight. Mitigation: keep a spare charge controller in inventory, configured and ready to install. Swap time target: under one hour.

- **Inverter failure.** Loss of AC power to all loads. Impact: total loss of usable power even if batteries are full. Detection: immediate -- all AC loads go dark. Mitigation: spare inverter in inventory. For critical loads, consider a secondary inverter on a separate circuit that can be activated manually during primary inverter failure.

- **Wiring degradation.** UV exposure, rodent damage, corrosion at terminals, and mechanical fatigue at junction points. Impact: increased resistance causes energy loss and heat buildup. Extreme cases cause fire. Detection: annual visual inspection and voltage drop testing. Mitigation: use appropriately rated cable, protect runs from animal access, torque all connections to specification.

- **Shading encroachment.** Vegetation growth, new structures, or shifting debris create shade on panels. Even partial shading on one cell can reduce an entire string's output dramatically. Detection: compare current output to historical records for the same month. Visual inspection. Mitigation: annual site survey to identify and remove shading sources.

- **Catastrophic weather event.** Hail, extreme wind, falling trees, or ice loading destroy panels. Impact: immediate and potentially total loss of generation capacity. Detection: obvious. Mitigation: geographic distribution of the array where possible, adequate structural mounting, and a replacement panel inventory sufficient to restore minimum operational capacity.

## 7. Recovery Procedures

1. **Single panel failure:** Isolate the failed panel. If the panel is part of a series string, the bypass diode should allow the string to continue operating at reduced voltage. Verify the string is still operating. Replace the panel at the next scheduled maintenance window. Use a panel of the same model and wattage. If the exact model is unavailable, a panel of similar specifications can be used in a parallel string but should not be mixed into an existing series string of different-specification panels.

2. **Charge controller failure:** Switch to the spare charge controller. Verify that charge parameters match the documented settings exactly. Monitor the first full charge cycle to confirm proper operation. Order a replacement spare immediately.

3. **Inverter failure:** Switch to the spare inverter. If no spare is available, critical DC loads (12V or 24V lighting, direct-DC-powered devices) can be powered directly from the battery bank while awaiting replacement. This is an emergency measure only and must be documented as such.

4. **Extended low-sun period:** If battery reserves fall below 50% state of charge due to extended cloud cover, implement the load reduction protocol: shut down non-essential systems, reduce computing loads to the minimum required for data integrity, and run the generator (D4-006) if available. Document the event and reassess the array sizing during the next annual review.

5. **Catastrophic array damage:** If more than 50% of the array is destroyed, immediately enter emergency power mode (D4-006). Assess the damage. Deploy replacement panels from inventory. If inventory is insufficient, operate the generator while sourcing replacements. This is a Tier 2 governance event and must be documented accordingly per GOV-001.

## 8. Evolution Path

- **Years 0-5:** The system is new. Panels are at peak production. Establish baseline measurements. Record monthly output data religiously. Calibrate expectations against reality. The load analysis will likely prove slightly wrong -- adjust the safety margin assessment accordingly. Replace any panels that exhibit early failure (infant mortality).

- **Years 5-10:** Degradation is measurable but manageable. Cumulative output loss of 3-5% from original. Focus on connection maintenance: retorque all terminals, inspect cable insulation, clean all junction boxes. Consider whether load growth requires array expansion.

- **Years 10-15:** Degradation reaches 6-10%. The charge controller and inverter are approaching mid-life. Begin sourcing replacement units. Check that replacement panels compatible with the existing array are still available. If the panel model has been discontinued, identify compatible alternatives and document the substitution plan.

- **Years 15-20:** Degradation reaches 10-15%. The original charge controller and inverter may need replacement. The mounting system should be inspected for structural fatigue, corrosion, and fastener integrity. Plan the next generation of the solar system: new panels, potentially new array configuration, potentially new battery chemistry (coordinate with D4-003).

- **Years 20-25:** Panels are reaching end of warranted life but may still produce 80-85% of rated output. The decision to replace the entire array or continue operating with degraded output is a Tier 2 governance decision. Factors: current output vs. current load, cost and availability of replacement panels, condition of mounting infrastructure.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators. Each entry should include the date, the author's identifier, and the context for the commentary.*

**2026-02-16 -- Founding Entry:**
Solar technology is one of the few areas where the institution's longevity requirement aligns comfortably with the technology's natural lifespan. A well-maintained solar panel genuinely lasts twenty-five years. The charge controller and inverter do not -- plan for two or three of each over the panel's lifetime. The most common mistake in off-grid solar design is undersizing: the operator calculates the minimum and builds that, leaving no margin for degradation, load growth, or bad weather. Size for the worst month, add a safety factor, and then add a little more. The cost of oversizing is a few extra panels. The cost of undersizing is a dark institution.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (Section 3.3: the off-grid mandate)
- OPS-001 -- Operations Philosophy (maintenance tempo, documentation-first principle)
- D4-001 -- Infrastructure Architecture Overview
- D4-003 -- Battery Systems: Selection, Management, and End-of-Life
- D4-005 -- Environmental Control Systems
- D4-006 -- Power Distribution and UPS Management
- IEC 61215 -- Terrestrial Photovoltaic Modules Design Qualification
- IEC 62109 -- Safety of Power Converters for Photovoltaic Power Systems

---

---

# D4-003 -- Battery Systems: Selection, Management, and End-of-Life

**Document ID:** D4-003
**Domain:** 4 -- Infrastructure & Power
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, OPS-001, D4-001, D4-002
**Depended Upon By:** D4-006 (Power Distribution), D4-005 (Environmental Controls). All articles involving energy storage, power backup, or autonomy duration.

---

## 1. Purpose

This article provides a complete engineering reference for the battery energy storage system of the holm.chat Documentation Institution. It covers battery chemistry selection, bank design and sizing, battery management system (BMS) configuration, capacity testing procedures, degradation tracking, replacement decision criteria, and end-of-life handling.

If the solar array (D4-002) is the institution's lungs, the battery bank is its heart. The solar array generates energy only during daylight hours. The battery bank stores that energy and releases it on demand -- overnight, during cloudy weather, and during any period when generation falls below consumption. The battery bank determines the institution's autonomy: how long it can operate without any generation at all.

Battery systems are the most maintenance-intensive and the most failure-prone component of an off-grid power system. They are also the most expensive component to replace over the institution's lifetime. A poorly managed battery bank will fail in three to five years. A well-managed bank will last ten to fifteen years or more, depending on chemistry. The difference is entirely in the operator's discipline: proper charging parameters, appropriate depth of discharge limits, regular capacity testing, and timely replacement when degradation crosses the threshold.

This article is written to make that discipline concrete, measurable, and sustainable.

## 2. Scope

This article covers:

- Battery chemistry comparison: lead-acid (flooded, AGM, gel) vs. lithium (LFP, NMC) for off-grid institutional use.
- Bank sizing: determining capacity requirements based on autonomy targets and load profiles.
- Bank configuration: series and parallel connections, cell matching, voltage architecture.
- Battery management system (BMS) configuration for lithium systems.
- Charge parameter specification for both chemistries.
- Capacity testing procedures: how to measure actual capacity vs. rated capacity.
- Degradation curves: expected capacity loss over time for each chemistry.
- Temperature effects on performance and lifespan.
- The replacement decision: when degradation justifies replacement.
- End-of-life handling: safe decommissioning, recycling, and disposal.
- The twenty-year battery lifecycle plan.

This article does not cover solar panel selection or charge controller configuration (see D4-002), internal power distribution or UPS systems (see D4-006), or the generator backup system (see D4-006). It covers the battery bank itself: from new cells to decommissioned cells.

## 3. Background

### 3.1 The Role of Storage in an Off-Grid Institution

The institution's mandate requires independence from the power grid (CON-001, Section 3.3). Solar panels generate power only when the sun shines. The gap between generation and consumption must be bridged by stored energy. In temperate climates, winter nights are long and winter days are short. The battery bank must store enough energy to power the institution through sixteen or more hours of darkness, and must have enough reserve to handle multiple consecutive cloudy days.

This is not a convenience feature. It is a survival requirement. When the battery bank is depleted and no generation is available, the institution shuts down. Data integrity, service availability, and the institution's ability to fulfill its mission all depend on the battery bank having sufficient capacity, being properly maintained, and being replaced before it fails.

### 3.2 The Chemistry Decision

The choice of battery chemistry is one of the most consequential decisions in the institution's physical design. It affects cost, lifespan, maintenance requirements, safety characteristics, temperature sensitivity, depth of discharge limits, charge parameters, physical space requirements, and end-of-life handling. The two practical options for off-grid institutional use are lead-acid and lithium iron phosphate (LFP). Each has genuine advantages. Neither is universally superior.

### 3.3 A Note on Battery Safety

All battery chemistries store significant energy and pose safety risks if mishandled. Lead-acid batteries produce hydrogen gas during charging, which is explosive in enclosed spaces. Lithium batteries can enter thermal runaway if overcharged, physically damaged, or operated outside their temperature specifications. Both chemistries contain materials that are hazardous if released into the environment. The safety precautions in this article are not optional. They are integral to every procedure.

## 4. System Model

### 4.1 Chemistry Comparison

**Lead-Acid (Flooded/Wet Cell):**
- Cycle life: 500-1,500 cycles at 50% depth of discharge (DoD), depending on quality.
- Usable capacity: 50% of rated capacity (deeper discharge dramatically shortens life).
- Maintenance: High. Requires periodic electrolyte level checks and distilled water addition (every 1-3 months). Terminal cleaning. Equalization charging.
- Temperature sensitivity: Capacity drops approximately 1% per degree C below 25 degrees C. Freezing can destroy cells if discharged.
- Self-discharge: 3-5% per month at 25 degrees C.
- Hydrogen off-gassing: Yes, requires ventilated enclosure.
- Cost: Low initial cost. Higher total cost of ownership due to shorter lifespan and maintenance.
- Expected lifespan in institutional use: 4-7 years with proper maintenance.
- Weight and size: Heavy and large for the capacity provided. Approximately 30-40 Wh per kg.

**Lead-Acid (AGM -- Absorbed Glass Mat):**
- Cycle life: 400-1,200 cycles at 50% DoD.
- Usable capacity: 50% of rated capacity.
- Maintenance: Low. Sealed, no water addition required.
- Temperature sensitivity: Similar to flooded lead-acid.
- Self-discharge: 1-3% per month at 25 degrees C.
- Hydrogen off-gassing: Minimal under normal operation. Valves vent if overcharged.
- Cost: Moderate. Less maintenance cost than flooded, but higher purchase price.
- Expected lifespan in institutional use: 4-6 years with proper maintenance.
- Weight and size: Similar to flooded lead-acid.

**Lithium Iron Phosphate (LiFePO4 / LFP):**
- Cycle life: 2,000-5,000 cycles at 80% DoD, depending on quality and operating conditions.
- Usable capacity: 80-90% of rated capacity.
- Maintenance: Very low. No water, no equalization. BMS handles cell balancing.
- Temperature sensitivity: Charging below 0 degrees C causes lithium plating and permanent damage. Discharging is safe to -20 degrees C but with reduced capacity. Optimal range: 15-35 degrees C.
- Self-discharge: Less than 2% per month.
- Hydrogen off-gassing: None under normal operation.
- Cost: High initial cost. Lower total cost of ownership due to longer lifespan and deeper usable capacity.
- Expected lifespan in institutional use: 8-15 years with proper management.
- Weight and size: Approximately 90-160 Wh per kg. Roughly one-third the weight and volume of lead-acid for equivalent usable capacity.

**The institutional recommendation:** LFP is the preferred chemistry for the institution's primary battery bank. The combination of longer cycle life, deeper usable capacity, lower maintenance, absence of hydrogen off-gassing, and lower total cost of ownership aligns with the institution's principles of longevity (ETH-001, Principle 4) and sustainable operations (OPS-001, Section 4.4). Lead-acid AGM is acceptable as a secondary or backup bank, and may be appropriate where budget constraints make LFP impractical for the initial build.

### 4.2 Bank Sizing

The battery bank must be sized to provide the institution's required autonomy -- the number of days the institution can operate without any solar generation.

**Step 1: Determine daily energy consumption.** Use the load analysis from D4-002, Section 4.1. Example: 5,000 Wh per day.

**Step 2: Determine autonomy target.** For temperate climates with occasional multi-day cloud cover, a minimum of three days of autonomy is recommended. For locations with extended low-sun seasons, four to five days may be appropriate. The autonomy target is a Tier 2 governance decision because it determines a fundamental infrastructure parameter.

**Step 3: Calculate total energy storage required.** Multiply daily consumption by autonomy days. Then divide by the usable capacity fraction for the chosen chemistry. For LFP at 80% usable: 5,000 Wh x 3 days / 0.80 = 18,750 Wh (18.75 kWh) of rated battery capacity. For lead-acid at 50% usable: 5,000 Wh x 3 days / 0.50 = 30,000 Wh (30 kWh) of rated capacity.

**Step 4: Apply a degradation margin.** Batteries lose capacity over their life. Size the bank so that even at end-of-life (typically 80% of original capacity for LFP, 60-70% for lead-acid), the bank still meets the autonomy target. For LFP: 18,750 Wh / 0.80 = 23,438 Wh. Round up to the nearest available battery configuration.

### 4.3 Bank Configuration

**Voltage architecture:** Common off-grid system voltages are 12V, 24V, and 48V. Higher voltage means lower current for the same power, which means smaller wire sizes and lower resistive losses. For systems above 2 kW continuous load, 48V is strongly recommended. For systems above 5 kW, 48V is essentially required.

**Series connections:** Cells or batteries wired in series increase voltage. Four 12V batteries in series produce 48V. All batteries in a series string must be identical: same manufacturer, same model, same age, same capacity. Mismatched batteries in series will cause the weakest battery to be overcharged during charging and over-discharged during discharging, accelerating its failure and potentially creating a safety hazard.

**Parallel connections:** Strings wired in parallel increase capacity. Two 48V strings in parallel double the amp-hour capacity. Parallel strings should also be matched, though the matching requirement is less critical than for series connections. Limit parallel strings to four or fewer to minimize circulating currents between strings.

### 4.4 Battery Management System (BMS) Configuration

LFP batteries require a BMS. The BMS monitors individual cell voltages, manages cell balancing, enforces charge and discharge limits, monitors temperature, and disconnects the battery in fault conditions.

Critical BMS parameters to configure and document:

- **High voltage cutoff (per cell):** 3.60-3.65V. Charging must stop when any cell reaches this voltage.
- **Low voltage cutoff (per cell):** 2.50-2.80V (manufacturer-specific). Discharging must stop when any cell reaches this voltage.
- **Charge current limit:** Per manufacturer specification, typically 0.5C maximum (half the amp-hour rating).
- **Discharge current limit:** Per manufacturer specification, typically 1C maximum continuous.
- **Low temperature charge cutoff:** 0 degrees C. Charging below freezing causes permanent lithium plating damage. This is non-negotiable.
- **High temperature cutoff:** 45-55 degrees C (manufacturer-specific). Both charging and discharging stop.
- **Cell balance threshold:** The voltage differential at which the BMS begins active or passive balancing. Typically 10-30mV.

Document all BMS parameters in the operational log. Verify them quarterly. If the BMS firmware is updatable, document the version and do not update without testing on a non-critical system first.

### 4.5 Capacity Testing Procedures

Capacity testing measures the battery bank's actual usable capacity compared to its rated capacity. This is the single most important diagnostic for tracking battery health.

**Full capacity test procedure (perform annually):**

1. Fully charge the battery bank. Allow it to rest at full charge for at least 2 hours.
2. Record the starting voltage of each cell (LFP) or each battery (lead-acid).
3. Discharge the bank at a controlled, constant rate approximately equal to the institution's average load. Record the discharge current.
4. Monitor cell/battery voltages throughout the discharge. Record voltage at regular intervals (every 15-30 minutes).
5. Stop the discharge when any cell reaches the low voltage cutoff (LFP) or when bank voltage reaches the 50% DoD voltage (lead-acid).
6. Calculate the total energy discharged: sum of (voltage x current x time interval) for each recording interval. This is the bank's measured capacity.
7. Compare measured capacity to the bank's rated capacity and to previous test results.
8. Record all data in the battery log.

**Interpreting results:** New LFP batteries should deliver 95-100% of rated capacity. New lead-acid batteries should deliver 90-100%. Capacity below 80% of rated (for LFP) or below 70% of rated (for lead-acid) indicates the battery is approaching end of life and replacement should be planned.

### 4.6 Degradation Curves

**LFP typical degradation:**
- Year 0-1: 1-3% capacity loss (initial settling).
- Year 1-5: 1-2% loss per year. Cumulative: 5-12% loss.
- Year 5-10: 1-3% loss per year, accelerating slightly. Cumulative: 10-25% loss.
- Year 10-15: 2-4% loss per year. Cumulative: 20-40% loss. Replacement zone typically reached at 10-15 years.

**Lead-acid (AGM) typical degradation:**
- Year 0-1: 3-8% capacity loss (initial and significant).
- Year 1-3: 3-5% loss per year. Cumulative: 10-20% loss.
- Year 3-5: 5-10% loss per year, accelerating. Cumulative: 20-40% loss.
- Year 5-7: Rapid decline. Replacement typically required at 4-7 years.

These curves assume proper charge parameters, appropriate depth of discharge, and operation within the recommended temperature range. Abuse -- chronic overcharging, deep discharging, temperature extremes, or extended storage in a discharged state -- accelerates degradation dramatically.

### 4.7 The Replacement Decision

Replace the battery bank when measured capacity falls below 80% of rated capacity (LFP) or 60% of rated capacity (lead-acid), or when the bank can no longer meet the institution's autonomy target, whichever comes first.

Do not wait for total failure. A failing battery bank gives warning through declining capacity test results. Heed those warnings. The cost of proactive replacement is a planned expense. The cost of reactive replacement after failure is potential data loss, unplanned downtime, and emergency procurement.

Replacement is a Tier 2 governance decision (GOV-001) because it involves significant expense and architectural impact.

### 4.8 End-of-Life Handling

**Lead-acid batteries:** Contain lead and sulfuric acid. Both are hazardous. Do not dispose of in general waste. Lead-acid batteries are the most recycled consumer product in the world -- recycling infrastructure exists in most regions. Transport in acid-resistant containers. Terminals must be taped or covered to prevent short circuits.

**LFP batteries:** Contain lithium, iron, and phosphate compounds. Less hazardous than other lithium chemistries but still require proper handling. Fully discharge before decommissioning (to below 2.5V per cell if possible). Do not puncture, crush, or incinerate. Recycling infrastructure for LFP is developing. Store decommissioned cells in a cool, dry, fireproof location until recycling is available.

Both chemistries: document the decommissioning in the operational log. Record the final capacity test results, the date of decommissioning, and the disposal or recycling method used.

## 5. Rules & Constraints

- **R-D4-003-01:** Battery charge parameters must exactly match the manufacturer's specifications. "Close enough" is not acceptable for charge voltages. Overcharging by even 0.1V per cell on LFP causes accelerated degradation. Undercharging causes sulfation in lead-acid.
- **R-D4-003-02:** Full capacity tests must be performed annually and the results recorded. A battery bank that has not been capacity-tested in the past 12 months has unknown health and cannot be relied upon for the stated autonomy.
- **R-D4-003-03:** LFP batteries must never be charged below 0 degrees C. If the battery environment cannot be guaranteed to remain above freezing, a low-temperature charge lockout must be implemented in hardware (BMS) or, as a secondary measure, in the charge controller.
- **R-D4-003-04:** Lead-acid batteries in enclosed spaces must have forced ventilation to prevent hydrogen accumulation. The ventilation system is a safety-critical component and must be tested quarterly.
- **R-D4-003-05:** All batteries in a series string must be matched: same manufacturer, model, capacity rating, and age. Mixing batteries in a series string is prohibited.
- **R-D4-003-06:** The replacement decision must be made proactively based on capacity test trends, not reactively after failure. When capacity drops below 85% of rated (LFP) or 70% of rated (lead-acid), begin procurement of replacements.
- **R-D4-003-07:** End-of-life batteries must be recycled through appropriate channels. Environmental contamination from improper battery disposal violates ETH-001, Principle 5 (Harm Reduction).

## 6. Failure Modes

- **Cell imbalance (LFP).** Individual cells drift out of balance, causing some cells to hit voltage limits before others. Effect: reduced usable capacity, potential overcharge of strongest cells. Detection: BMS cell voltage monitoring shows increasing spread between highest and lowest cell. Mitigation: verify BMS balancing is functioning. Perform manual top-balance if necessary (fully charge, hold at absorption voltage until all cells equalize, may require 24-48 hours).

- **Sulfation (lead-acid).** Chronic undercharging causes lead sulfate crystals to harden on the plates, permanently reducing capacity. Effect: accelerated capacity loss, increased internal resistance. Detection: capacity test shows declining capacity; batteries take longer to reach full charge. Mitigation: ensure charge parameters are correct. Perform equalization charges per manufacturer schedule. Severe sulfation is irreversible.

- **Thermal runaway (LFP -- rare but catastrophic).** Internal short circuit causes uncontrolled temperature rise. Effect: fire, toxic fumes, potential explosion. Cause: physical damage, manufacturing defect, or BMS failure allowing extreme overcharge. Mitigation: functioning BMS, proper fusing on all battery connections, fire-resistant battery enclosure, ventilation to exterior. Recovery: evacuate the area, allow the event to complete, do not attempt to extinguish with water. Replace the entire affected string.

- **BMS failure.** The BMS stops monitoring or controlling the battery. Effect: potential overcharge, over-discharge, or operation outside temperature limits -- all of which cause permanent damage. Detection: monitor BMS status indicators. If the BMS becomes unresponsive, disconnect the battery bank immediately. Mitigation: keep a spare BMS. Configure the charge controller's own voltage limits as a secondary protection layer that would stop charging even if the BMS fails.

- **Connection corrosion.** Battery terminals corrode over time, increasing resistance. Effect: voltage drop at terminals, heat generation, reduced effective capacity. Detection: visual inspection, voltage measurement at terminals vs. at battery posts. Mitigation: annual terminal cleaning, application of anti-corrosion compound, torque verification.

## 7. Recovery Procedures

1. **Single cell failure in LFP bank:** If one cell in a series string fails (short-circuit or open-circuit), the entire string must be taken offline. If parallel strings exist, the remaining strings continue operation at reduced capacity. Replace the failed cell with a matched cell. If a matched cell is not available, replace the entire string. A mismatched cell in a series string is a long-term degradation risk.

2. **BMS failure:** Disconnect the battery bank from both the charge source and the loads. Connect the spare BMS. Verify all parameters match the documented settings. Reconnect the bank and monitor the first charge/discharge cycle closely.

3. **Capacity below replacement threshold:** Begin procurement of replacement batteries. Continue operating the existing bank at reduced autonomy. Calculate the new effective autonomy based on measured capacity. Adjust the load reduction protocol trigger point accordingly. Do not wait until the bank fails completely.

4. **Hydrogen accumulation alarm (lead-acid):** Immediately ventilate the battery enclosure. Do not operate electrical switches or create sparks in the area. Identify the cause: failed ventilation fan, overcharging, or excess equalization. Repair the ventilation system before resuming charging. This is a safety-critical event.

5. **Deep discharge event:** If the bank has been deeply discharged (below the recommended minimum -- below 2.5V per cell for LFP, below 10.5V per 12V battery for lead-acid), immediately recharge using the gentlest charge rate available (0.1C or lower). Monitor cell voltages closely. Some cells may not recover. Perform a full capacity test after recovery charging. Document the event and the resulting capacity loss.

## 8. Evolution Path

- **Years 0-5:** The initial battery bank is new. Establish baseline capacity measurements. Perform capacity tests every six months in the first two years to establish the degradation curve early. Verify charge parameters are producing the expected results.

- **Years 5-10:** For LFP: the bank should still be performing well, with 80-95% of original capacity. For lead-acid: the first replacement cycle will occur during this period. Begin procurement of replacements when capacity crosses 85% (LFP) or 70% (lead-acid).

- **Years 10-15:** For LFP: the bank is entering the replacement zone. Capacity tests become critical. For lead-acid: you may be on the second or third bank. Evaluate whether switching to LFP is justified for the next replacement.

- **Years 15-20:** Battery chemistry technology will have evolved. New chemistries may offer superior performance. Evaluate the options against the institution's requirements: cycle life, maintenance burden, safety profile, temperature range, and cost of ownership. The replacement decision at this point should be informed by fifteen years of operational data. Use that data.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Batteries are the component of the off-grid power system that most rewards diligent management and most punishes neglect. A solar panel will tolerate years of indifference and still produce 85% of its rated output. A battery bank will tolerate perhaps six months of incorrect charge parameters before the damage is irreversible. Know your charge voltages. Test your capacity. Watch the degradation curve. When the curve tells you the bank is aging, believe it and act. The most expensive battery is the one that fails before its replacement arrives.

The choice between lead-acid and LFP is real. Lead-acid is cheaper at the point of purchase and available everywhere. LFP is cheaper over its lifetime and requires less of the operator's attention. For an institution designed to last decades, operated by a single person, the lower total cost of ownership and lower maintenance burden of LFP is worth the higher initial investment. But lead-acid is honest, well-understood technology that has kept off-grid installations running for decades. Either chemistry works if managed properly. Neither works if neglected.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 4: Longevity Over Novelty; Principle 5: Harm Reduction -- disposal)
- CON-001 -- The Founding Mandate (off-grid mandate)
- OPS-001 -- Operations Philosophy (maintenance tempo, sustainability requirement)
- D4-001 -- Infrastructure Architecture Overview
- D4-002 -- Solar Power System Design and Maintenance
- D4-005 -- Environmental Control Systems (temperature management for batteries)
- D4-006 -- Power Distribution and UPS Management
- IEC 62619 -- Safety Requirements for Secondary Lithium Cells and Batteries
- IEEE 1188 -- Recommended Practice for Maintenance, Testing, and Replacement of VRLA Batteries

---

---

# D4-004 -- Network Hardware Lifecycle Management

**Document ID:** D4-004
**Domain:** 4 -- Infrastructure & Power
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, OPS-001, SEC-001, D4-001
**Depended Upon By:** D4-005 (Environmental Controls -- heat loads from network hardware), D5-001 (Platform systems -- network dependencies). All articles involving internal network communication, data transfer, or service connectivity.

---

## 1. Purpose

This article provides a complete lifecycle management reference for the network hardware of the holm.chat Documentation Institution. It covers switches, routers, access points, cabling infrastructure, and all associated components from initial selection through deployment, operation, maintenance, replacement, and decommissioning.

The institution is air-gapped. There is no connection to the public internet. But within the air gap, there is a network -- an internal network that connects the institution's servers, storage systems, workstations, and management interfaces. This internal network is the circulatory system of the institution. When it fails, systems that are individually healthy become collectively useless: servers cannot reach storage, workstations cannot reach servers, and the operator cannot manage anything.

Network hardware occupies an awkward position in most infrastructure planning. It is less dramatic than servers, less visible than storage, and less obviously critical than power systems. It is also the first thing blamed when something goes wrong and the last thing proactively maintained. This article exists to correct that neglect. Network hardware has a finite lifespan, requires specific environmental conditions, degrades in ways that are not always obvious, and depends on a supply chain that can make specific models unavailable with little warning. Managing it requires planning, discipline, and a spare parts inventory.

## 2. Scope

This article covers:

- Network architecture overview for the air-gapped institution.
- Hardware selection criteria: switches, routers, access points, and cabling.
- Lifecycle stages: procurement, deployment, operation, maintenance, replacement, decommissioning.
- The network hardware audit: a structured assessment of all network components.
- Performance monitoring and degradation detection.
- Spare parts inventory management.
- Procurement strategy for long-lived equipment and discontinued models.
- Cabling infrastructure: types, standards, testing, and replacement.
- Documentation requirements for the network.

This article does not cover network protocol configuration, firewall rules, VLAN architecture, or software-level network security (those belong in Domain 3 and Domain 5 articles). It covers the physical hardware and its lifecycle.

## 3. Background

### 3.1 The Network in an Air-Gapped Institution

The institution's network is entirely internal. There is no WAN link, no internet gateway, no DNS resolution of external names. The network exists to connect the institution's internal systems to each other and to the operator's workstation. This simplifies some aspects of network management (no routing table complexity, no BGP peering, no firewall rules for external traffic) and complicates others (no cloud-based monitoring tools, no automatic firmware updates, no vendor support portals).

The air gap has a specific implication for network hardware: once a device is deployed inside the air gap, it cannot receive firmware updates from the manufacturer without a deliberate, quarantined import process (see SEC-001, Domain 18). This means the firmware version deployed at installation may be the firmware version running for the device's entire operational life. Select hardware with mature, stable firmware. Avoid hardware that requires frequent updates to remain functional.

### 3.2 The Lifecycle Problem

Network hardware has a typical operational life of seven to twelve years, depending on the environment, workload, and component quality. But the products themselves are often discontinued within three to five years of release. This creates a lifecycle problem: the hardware you deploy today will be out of production long before it reaches end of life. Replacement parts, compatible expansion modules, and exact-model replacements will become unavailable.

The institution must plan for this reality. The procurement strategy described in this article is designed to ensure that when a network device fails at year eight, a replacement is already on hand -- purchased years earlier while the model was still available, or a compatible successor identified and validated in advance.

### 3.3 Why Not Consumer Equipment

Consumer-grade network equipment (home routers, consumer Wi-Fi access points, unmanaged switches) is designed for a three to five year lifecycle, offers limited or no management capabilities, and frequently depends on cloud services for configuration and monitoring. Enterprise or professional-grade equipment is designed for longer operation, offers full local management (SNMP, CLI, local web interfaces), supports advanced features like VLAN segmentation and port mirroring, and does not phone home. The additional cost of professional-grade equipment is justified by the institution's longevity requirement (ETH-001, Principle 4).

## 4. System Model

### 4.1 Network Architecture Overview

The institution's internal network follows a simple, resilient architecture:

**Core switch:** A managed switch that serves as the central connection point for all network traffic. This is the single most critical network device. It should be a high-quality, enterprise-grade managed switch with sufficient ports for all current connections plus 30% growth capacity. Gigabit Ethernet minimum. 10 Gigabit for connections to storage arrays if the storage protocol requires it.

**Distribution switches (if needed):** For installations where equipment is spread across multiple rooms or racks, secondary managed switches connect to the core switch via trunk links. Each distribution switch serves the devices in its local area.

**Management network (optional but recommended):** A physically separate or VLAN-segregated network for out-of-band management of servers and infrastructure devices (IPMI, iLO, DRAC, switch management interfaces). This allows the operator to manage devices even when the primary network is impaired.

**Wireless access point (if needed):** For operator workstation mobility within the facility. Must support WPA3 or equivalent. Must be locally managed, not cloud-managed. Wireless is a convenience, not a requirement. All critical connections should be wired.

**Cabling infrastructure:** Category 6A or better for all permanent installations. Supports 10 Gigabit Ethernet to 55 meters, providing headroom for future bandwidth upgrades. Fiber optic for runs exceeding 50 meters or for connections between buildings.

### 4.2 Hardware Selection Criteria

When selecting network hardware, evaluate against these criteria in order of priority:

1. **Local manageability:** The device must be fully configurable and monitorable through a local interface (serial console, SSH, local web interface, SNMP). Cloud-dependent management is a disqualifier. The institution is air-gapped; cloud management is not merely inconvenient, it is architecturally impossible.

2. **Firmware stability:** The device should have mature, stable firmware with a track record of reliable operation. Avoid newly released product lines with immature firmware. Check community forums and independent reviews for reports of firmware issues.

3. **Build quality and environmental ratings:** The device should have a fanless design if possible (fans are the most common failure point in network equipment) or easily replaceable fans. Operating temperature range should exceed the environmental conditions in the equipment space (see D4-005).

4. **Port density and type:** Sufficient ports for current needs plus 30% growth. Appropriate port speeds for the traffic profile.

5. **Power consumption:** Lower power consumption reduces the load on the solar and battery systems (D4-002, D4-003). Compare power consumption across candidates.

6. **Replaceability and parts availability:** Prefer manufacturers with a history of long product lifecycles and wide distribution. Avoid boutique or niche brands that may exit the market.

### 4.3 The Network Hardware Audit

The network hardware audit is a structured assessment of every network component in the institution. It is performed annually as part of the annual operations cycle (OPS-001, Section 4.1) and whenever a network issue suggests hardware degradation.

The audit inspects each device for:

- **Physical condition:** Visual inspection for damage, discoloration (heat damage), dust accumulation, LED status indicators, fan operation (if applicable), cable connection integrity.
- **Error counters:** Check interface error counters (CRC errors, frame errors, collisions, drops, discards). Non-zero error counters indicate cabling problems, hardware degradation, or electromagnetic interference. Record error counts and compare to previous audit.
- **Port health:** Test each active port for proper link negotiation (correct speed, correct duplex). Test unused ports to verify they are still functional.
- **Environmental readings:** Record device temperature (from SNMP or management interface if available). Compare to historical readings.
- **Uptime and stability:** Record uptime. Unexpected reboots indicate hardware instability or power issues.
- **Firmware version:** Record current firmware version. Note whether a more recent firmware is available and whether the update has been evaluated for import.
- **Cable testing:** For permanent cable installations, test continuity, wire mapping, and link speed negotiation. For critical runs, test with a cable certifier annually.

Document the audit results in the infrastructure log. Flag any findings that require action. Track trends across audits.

### 4.4 Spare Parts Inventory

The spare parts inventory is the institution's insurance against hardware failure. For network hardware, the inventory should include:

- **One spare core switch** (identical model or validated compatible model). This is the highest-priority spare. A core switch failure with no spare available is a total network outage.
- **One spare of each distribution switch model** (if distribution switches are used).
- **Spare SFP/SFP+ modules** for all fiber connections (minimum two of each type in use).
- **Spare patch cables** in each category and length in use (minimum five of each).
- **Spare power supplies** for any device with a replaceable power supply module.
- **One spare wireless access point** (if wireless is deployed).
- **Spare fans** for any device with replaceable fan modules.
- **A serial console cable** compatible with the management console port of each device type.

Store spares in their original packaging in a clean, temperature-controlled environment. Test spares annually: power them on, verify basic functionality, and return to storage. An untested spare is a hope, not a plan.

### 4.5 Procurement Strategy for Discontinued Equipment

When a network device model is discontinued by the manufacturer, the institution has three options:

**Option 1: Stockpile.** Purchase additional units of the discontinued model before they disappear from the market. This is the simplest approach and is recommended when the device is performing well and compatible replacements are uncertain. Purchase enough units to cover the expected remaining operational life of the deployed units.

**Option 2: Identify a successor.** Research the manufacturer's replacement product or competing products with equivalent capabilities. Validate the successor in a test environment before committing. Document any configuration differences between the old and new models.

**Option 3: Redesign.** If no direct successor exists, redesign the affected portion of the network to use available equipment. This is the most disruptive option and should be treated as a Tier 2 governance decision.

The procurement decision should be made as soon as end-of-sale or end-of-life is announced -- not when the last deployed unit fails. Track manufacturer lifecycle announcements for all deployed models.

### 4.6 Cabling Lifecycle

Cabling is the most overlooked and longest-lived component of the network. Properly installed structured cabling lasts twenty to thirty years or more. Poorly installed cabling causes intermittent, difficult-to-diagnose problems from day one.

**Installation standards:** All permanent cabling should be installed to TIA-568 standards. Cable runs should maintain bend radius minimums, avoid sources of electromagnetic interference (power cables, motors, fluorescent lights), and be properly terminated and labeled at both ends.

**Labeling:** Every cable run must be labeled with a unique identifier at both ends. The identifier must map to a cable schedule document that records the source, destination, cable type, length, installation date, and most recent test result.

**Testing:** All new cable installations must be tested with a cable certifier before being placed in service. Annual testing of permanent runs ensures that degradation (from environmental factors, physical stress, or connector wear) is detected before it causes problems.

**Replacement criteria:** Replace a cable run when it fails certification testing, when error counters on the connected ports show persistent errors traceable to the cable, or when the cable jacket shows visible damage (cracking, discoloration, deformation).

## 5. Rules & Constraints

- **R-D4-004-01:** All network hardware must be locally manageable. Cloud-dependent management is a disqualifier. This derives directly from the air-gap mandate (SEC-001, R-SEC-01) and the sovereignty principle (ETH-001, Principle 1).
- **R-D4-004-02:** A spare core switch must be maintained in inventory at all times. Operating without a core switch spare is operating without a safety net for the institution's internal communication.
- **R-D4-004-03:** The network hardware audit must be performed at least annually. Results must be documented and trends tracked across audits.
- **R-D4-004-04:** All permanent cable installations must be tested at installation and annually thereafter. Untested cables are potential failure points.
- **R-D4-004-05:** When a deployed network device model reaches end-of-sale announcement, the procurement response (stockpile, successor, or redesign) must be decided and documented within 90 days.
- **R-D4-004-06:** All spare parts must be tested annually. Dead spares provide false confidence.
- **R-D4-004-07:** Every cable run must be labeled and documented in the cable schedule. Unlabeled cables are untraceable, and untraceable cables are unmaintainable.

## 6. Failure Modes

- **Core switch failure.** Total internal network outage. All inter-system communication stops. Servers cannot reach storage. Operator cannot manage systems remotely. Impact: equivalent to total institution outage for any networked service. Detection: immediate -- all network-dependent services fail. Mitigation: spare core switch with documented configuration. Target recovery time: under two hours to physically swap and restore configuration.

- **Port failure (individual).** A single port on a switch stops functioning or develops errors. Impact: one device loses connectivity. Detection: link down alarm or error counter escalation on the affected port. Mitigation: move the affected connection to a spare port. Mark the failed port as out of service. If port failures become frequent on a device, the switch may be approaching end of life.

- **Fan failure.** A cooling fan in a managed switch fails. Impact: device temperature rises. If the device has thermal protection, it will throttle performance or shut down. If it does not, it will operate at elevated temperature with accelerated degradation until it fails completely. Detection: fan RPM monitoring via SNMP, audible change in noise, visual inspection. Mitigation: replace the fan immediately. If the fan is not field-replaceable, replace the entire device.

- **Power supply failure (in devices with single power supply).** Total device outage. Impact: depends on device role -- core switch failure is catastrophic, edge switch failure is localized. Detection: immediate -- device goes dark. Mitigation: devices with dual power supply capability should use both power supplies connected to separate power sources. For single-power-supply devices, keep a spare device ready.

- **Cable failure.** Intermittent connectivity, high error rates, or complete link failure. Impact: ranges from performance degradation to total link loss. Detection: error counter monitoring, link flapping, cable certification failure. Mitigation: replace the cable. This is one of the simplest and most common network fixes.

- **Firmware bug.** Undiscovered firmware bug manifests after extended operation. Impact: unpredictable -- could range from minor logging errors to device crashes. Detection: unexpected behavior, crashes, log analysis. Mitigation: if a firmware update is available that addresses the bug, evaluate it for import through the quarantine process. If no update is available, work around the bug and document the workaround.

## 7. Recovery Procedures

1. **Core switch failure:** Power off the failed switch. Retrieve the spare from inventory. Install the spare in the same location with the same cabling. Apply the documented configuration (from the most recent configuration backup). Verify all links come up correctly. Test connectivity between all critical systems. Record the event in the infrastructure log. Order a replacement spare immediately.

2. **Distributed switch failure:** If the failed switch serves non-critical connections, those connections are offline until the switch is replaced. If it serves critical connections, temporarily reroute critical connections to the core switch (if ports are available) while the replacement is deployed.

3. **Complete network outage of unknown cause:** Systematic diagnosis: verify power to all network devices. Check for physical damage or disconnection. Start at the core switch and work outward. Test the core switch with a single directly connected device. Add connections one at a time to isolate the failure point.

4. **Performance degradation:** Check error counters on all devices. Identify the link or device with the highest error rate. Test the cabling on that link. Replace cables first (cheapest, most common cause). If cables test clean, investigate the switch port and the connected device's NIC.

5. **Procurement emergency (device fails with no spare):** Assess which connections can be temporarily consolidated. A failed 24-port switch can sometimes be temporarily replaced by combining available spare ports on other switches with temporary cabling. This is an emergency measure. Document it. Procure a permanent replacement with urgency.

## 8. Evolution Path

- **Years 0-5:** Network hardware is new and stable. Establish baseline error counter readings and temperature profiles. Build the spare parts inventory to its full complement. Document the complete cable schedule.

- **Years 5-10:** Hardware approaches mid-life. Monitor for increasing error rates. Some manufacturers will have discontinued the deployed models. Execute the procurement strategy (Section 4.5). Begin evaluating successor models.

- **Years 10-15:** First-generation network hardware reaches end of life. Plan and execute replacements. This is an opportunity to upgrade bandwidth (e.g., from Gigabit to 10 Gigabit) if the institution's traffic patterns justify it. Replacement is a Tier 2 governance decision.

- **Years 15-20:** Second-generation hardware is in service. The cabling infrastructure, if properly installed, should still be sound. Verify with certification testing. Network standards will have evolved; evaluate whether the current architecture still serves the institution's needs.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Network hardware is the infrastructure that everyone forgets about until it fails. Switches sit in racks, blinking quietly, doing their job for years without complaint. This silence breeds complacency. The annual network audit exists specifically to combat that complacency: to force the operator to actually look at the equipment, read the error counters, check the temperatures, and verify that the spares are still functional. The single most important network investment is a spare core switch. Everything else can be worked around temporarily. A core switch failure with no spare is a full institution outage with no defined recovery path except "wait for a replacement to arrive," which, for an air-gapped institution with no internet ordering capability, may mean days or weeks. Buy the spare. Test the spare. Keep the spare ready.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 1: Sovereignty; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (air-gap mandate)
- SEC-001 -- Threat Model and Security Philosophy (air-gap enforcement, supply chain threats)
- OPS-001 -- Operations Philosophy (maintenance tempo, complexity budget)
- D4-001 -- Infrastructure Architecture Overview
- D4-005 -- Environmental Control Systems (operating temperature for network hardware)
- D4-006 -- Power Distribution and UPS Management (power requirements for network hardware)
- TIA-568 -- Commercial Building Telecommunications Cabling Standard
- ISO/IEC 11801 -- Information Technology -- Generic Cabling for Customer Premises

---

---

# D4-005 -- Environmental Control Systems

**Document ID:** D4-005
**Domain:** 4 -- Infrastructure & Power
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, OPS-001, D4-001, D4-002, D4-003
**Depended Upon By:** D4-004 (Network Hardware -- environmental requirements), D4-006 (Power Distribution -- power for environmental controls). All articles involving equipment placement, server room design, or hardware longevity.

---

## 1. Purpose

This article provides a complete reference for the environmental control systems of the holm.chat Documentation Institution. It covers temperature management, humidity control, dust mitigation, and the monitoring systems that tie them together. It defines the environmental envelope within which the institution's computing equipment can operate safely and for maximum longevity, and it specifies the systems, procedures, and contingencies required to maintain that envelope.

Computing equipment is physical. It occupies space. It generates heat. It is damaged by moisture, dust, and temperature extremes. These facts are easy to overlook when thinking about digital infrastructure, but they are as fundamental as the electrical power supply. A server operating at 45 degrees C in a dusty, humid room will fail years before an identical server operating at 22 degrees C in a clean, dry room. Environmental control is not a luxury. It is a direct determinant of hardware lifespan, and hardware lifespan directly determines the institution's replacement costs and operational continuity.

The challenge for an off-grid institution is that environmental control consumes power. Air conditioning, dehumidification, and active filtration all draw from the same solar-charged battery bank that powers the computing equipment itself. Every watt spent on cooling is a watt not available for computing. The design philosophy of this article reflects that tension: prefer passive solutions where possible, use active solutions only where necessary, and size everything to the actual requirements rather than to datacenter standards that assume unlimited grid power.

## 2. Scope

This article covers:

- Recommended environmental parameters for the institution's computing equipment.
- Heat load calculation: how much heat the equipment generates.
- Passive cooling strategies: ventilation, thermal mass, orientation, insulation.
- Active cooling strategies: when passive is insufficient and mechanical cooling is required.
- Humidity monitoring and control: acceptable ranges, condensation prevention, dehumidification.
- Dust management: filtration, sealing, cleaning schedules.
- Environmental monitoring systems: sensors, logging, alerting.
- The environmental control audit: structured assessment of the equipment environment.
- Contingency procedures: what happens when environmental controls fail.

This article does not cover the power systems that supply electricity to environmental controls (see D4-002, D4-003, D4-006) or the specific hardware being protected (see D4-004, Domain 5). It covers the environment itself and the systems that maintain it.

## 3. Background

### 3.1 Why Environment Matters

The relationship between operating environment and hardware lifespan is well-documented and quantifiable. The Arrhenius equation applied to electronics predicts that for every 10 degrees C increase in operating temperature above the recommended range, component life expectancy is approximately halved. At 35 degrees C ambient, a hard drive rated for 5 years at 25 degrees C may last only 2.5 years. At 45 degrees C, perhaps 1.5 years.

Humidity compounds the damage. High humidity (above 60% relative humidity) promotes corrosion of circuit board traces, connector contacts, and solder joints. It also creates conditions for condensation when temperatures drop, and condensation on live electronics causes short circuits. Low humidity (below 20% relative humidity) promotes electrostatic discharge, which can damage sensitive components instantly and invisibly.

Dust accumulates on heat sinks, fan blades, and air intake filters, reducing cooling efficiency. A dust-clogged heat sink can raise a processor's operating temperature by 15-20 degrees C, pushing a comfortably cooled chip into the thermal throttling zone or beyond.

For a conventional datacenter, the solution is industrial HVAC, precision cooling, and dedicated maintenance staff. For an off-grid, single-operator institution, those solutions are neither available nor appropriate. The institution needs solutions that work within its power budget, can be maintained by one person, and will last as long as the equipment they protect.

### 3.2 The Passive-First Philosophy

The institution's environmental control strategy follows a passive-first philosophy: design the physical space to minimize the need for active mechanical systems. Passive solutions consume no power, have no moving parts, and require minimal maintenance. They include building orientation, insulation, thermal mass, natural ventilation paths, and equipment placement.

Active systems -- air conditioners, dehumidifiers, powered ventilation fans -- are deployed only where passive solutions are insufficient to maintain the required environmental envelope. When active systems are deployed, they must be sized for the actual heat load, not for worst-case estimates that lead to oversized, power-hungry installations.

## 4. System Model

### 4.1 Environmental Parameters

The following parameters define the acceptable environmental envelope for the institution's computing equipment. These are based on ASHRAE TC 9.9 recommendations for data processing environments, adapted for the institution's context.

**Temperature:**
- Recommended operating range: 18-27 degrees C (64-80 degrees F).
- Allowable range (short-term, not to exceed 96 hours): 15-32 degrees C (59-90 degrees F).
- Maximum rate of change: 5 degrees C per hour (rapid temperature changes cause thermal stress on solder joints and connections).
- Storage (for powered-off equipment): 5-40 degrees C.

**Relative Humidity:**
- Recommended operating range: 30-55% RH.
- Allowable range (short-term): 20-60% RH.
- Dew point: must remain below 15 degrees C at all times. Condensation on equipment is a critical failure condition.

**Particulate contamination:**
- Target: ISO 14644-1 Class 8 or better (equivalent to a clean office environment).
- Visible dust accumulation on equipment surfaces should not be detectable between monthly cleaning cycles.

### 4.2 Heat Load Calculation

The equipment's heat load determines the cooling requirement. In a computing environment, virtually all electrical energy consumed is converted to heat. The heat load in watts is therefore approximately equal to the total power consumption of all equipment in the space.

**Step 1:** Sum the power consumption of all equipment in the computing space: servers, storage arrays, network switches, UPS systems, monitors, and any other electrical loads. Use measured values from the power audit (D4-006), not nameplate ratings.

**Step 2:** Add the heat contribution from lighting and from the operator's presence (approximately 100W per person).

**Step 3:** Add solar heat gain if the space has windows or uninsulated exterior walls exposed to direct sunlight. Solar gain can be significant -- up to 200-600 watts per square meter of sun-exposed window.

**Step 4:** The result is the total heat load in watts that must be removed from the space to maintain temperature. In well-insulated, below-grade spaces with no solar gain, the equipment heat load dominates. In above-grade spaces with windows, solar gain may exceed the equipment load in summer.

### 4.3 Passive Cooling Strategies

**Building orientation and placement:** If the equipment room's location is a choice, place it on the north side of the building (Northern Hemisphere) or south side (Southern Hemisphere) to minimize solar gain. Below-grade (basement) locations offer natural temperature stability -- ground temperature at depth varies far less than air temperature.

**Insulation:** Insulate the equipment room to reduce heat transfer from the external environment. In hot climates, insulation keeps external heat out. In cold climates, insulation retains the equipment's waste heat, reducing or eliminating the need for dedicated heating.

**Thermal mass:** Dense materials (concrete, stone, earth) absorb heat during warm periods and release it during cool periods, dampening temperature swings. A below-grade room with concrete walls provides excellent thermal mass.

**Natural ventilation:** When outside air temperature is lower than the equipment room temperature and humidity is within acceptable range, ventilation can remove heat without mechanical cooling. Design ventilation paths so that cool air enters at the bottom of the space and warm air exits at the top (hot air rises). Ventilation openings must be filtered to prevent dust and pest ingress. Ventilation must be closable for periods when outside conditions are not suitable (too hot, too humid, too dusty).

**Equipment arrangement:** Arrange equipment so that hot exhaust air is directed toward the room's exhaust path, not toward the intake of adjacent equipment. In rack-mounted installations, maintain consistent front-to-back airflow and fill empty rack spaces with blanking panels to prevent hot exhaust from recirculating to the front.

### 4.4 Active Cooling

When passive strategies cannot maintain the temperature within the recommended range, active cooling is required. For the off-grid institution, the options are:

**Split-system air conditioner (mini-split):** The most common and most practical option for small equipment rooms. A high-efficiency inverter-type mini-split consumes 300-1,500 watts depending on the heat load. Select a unit rated for the calculated heat load with a seasonal energy efficiency ratio (SEER) of 20 or higher to minimize power consumption. Ensure the unit can operate in cooling mode at the expected range of outdoor temperatures.

**Evaporative cooling:** In hot, dry climates (relative humidity consistently below 30%), evaporative cooling can reduce temperature significantly at very low energy cost. Not suitable for humid climates. Increases humidity, which must be monitored.

**Dedicated computer room air conditioner (CRAC):** Overkill for most institutional installations. Only consider if the heat load exceeds 5 kW and no other solution is adequate.

The power consumption of active cooling must be included in the load analysis for the solar and battery systems (D4-002, D4-003). This is critical: the cooling system may represent 20-40% of the total institutional power load in hot climates.

### 4.5 Humidity Control

**High humidity (above 55% RH):** Deploy a dehumidifier. For small spaces, a desiccant dehumidifier may be adequate and is simpler to maintain than a compressor-based unit. For larger spaces or persistent high humidity, a compressor dehumidifier is more efficient. Drain collected water automatically if possible; manual draining creates a maintenance task that is easily forgotten.

**Low humidity (below 30% RH):** Common in cold, dry climates and in spaces with significant heating. Low humidity increases electrostatic discharge risk. Humidification is the solution but adds complexity and a water source to the equipment room, which carries its own risks. In practice, maintaining relative humidity above 20% is sufficient to prevent most ESD damage. Wrist grounding straps and ESD-safe work practices may be more practical than humidification equipment for a small installation.

**Condensation prevention:** Condensation occurs when surface temperatures fall below the air's dew point. The most common cause in equipment rooms is cold equipment being brought from outdoors into a warm, humid room, or temperature drops in a humid room when cooling shuts off overnight. Mitigation: allow cold equipment to acclimate to room temperature before powering on (at least 2 hours). Maintain humidity within the recommended range. Avoid sudden temperature changes.

### 4.6 Dust Management

**Source control:** Seal the equipment room against dust ingress. Seal cable penetrations, pipe entries, and gaps around doors with appropriate materials. Use a positive-pressure ventilation design (slightly more air entering through filtered intakes than leaving) to prevent unfiltered air from being drawn in through gaps.

**Filtration:** All air entering the equipment room -- whether through ventilation openings, air conditioner intakes, or doorways -- should pass through appropriate filters. MERV 8 or higher for general dust. MERV 13 for environments with significant airborne particulate (dusty climates, agricultural settings, construction activity).

**Cleaning schedule:**
- Monthly: Wipe external surfaces of all equipment. Inspect air intake filters on all devices with internal fans. Clean or replace filters as needed.
- Quarterly: Vacuum the floor and all horizontal surfaces. Inspect cable runs for dust accumulation. Clean air handling filters.
- Annually: Open each device (if safe and warranty-appropriate) and inspect internal components for dust accumulation. Use compressed air (filtered, moisture-free) to clean heat sinks, fan blades, and circuit boards. This is part of the annual infrastructure review (OPS-001).

### 4.7 Environmental Monitoring

Deploy sensors to continuously monitor the equipment room environment:

**Minimum sensor deployment:**
- Temperature sensor(s): at least one per equipment rack or cluster of equipment. Place at the intake side of the equipment where air enters, not at the exhaust where it is hottest. For critical equipment, add a second sensor at the exhaust side to monitor the temperature differential.
- Humidity sensor(s): at least one per room, placed at a representative location away from direct air conditioner or dehumidifier output.
- Water leak sensor(s): on the floor beneath any water-carrying equipment (dehumidifier, air conditioner condensate line) and in any area where water could enter the room from external sources.

**Data logging:** All sensor data should be logged continuously at intervals of no more than 5 minutes. Retain log data for at least two years to enable trend analysis and seasonal comparison. Plot temperature and humidity trends monthly.

**Alerting:** Configure alerts for conditions outside the allowable range. Alerts must be detectable by the operator even when not actively monitoring. In an air-gapped institution, this means local audible or visual alarms (alarm buzzer, indicator light), not email or text message notifications.

### 4.8 The Environmental Control Audit

The environmental control audit is performed semi-annually (once at the peak of the hot season, once at the peak of the cold season) to verify that the environmental control systems are maintaining the required parameters under the most demanding conditions.

The audit inspects:

- All temperature and humidity sensor readings against the required parameters.
- All sensor calibrations (compare readings to a known-good reference instrument).
- Active cooling system performance: is it maintaining setpoint, or is it running continuously without achieving setpoint (indicating insufficient capacity)?
- Passive ventilation paths: are they clear, are filters clean, are dampers operating?
- Dust levels: visual inspection and surface wipe test.
- Water leak sensors: test each sensor by applying a small amount of water.
- Environmental log data trends: is the average temperature rising over time? Is humidity control becoming less effective? Trends indicate developing problems before they become acute.

Document the audit results in the infrastructure log. Compare to previous audits to identify trends.

## 5. Rules & Constraints

- **R-D4-005-01:** The equipment room temperature must remain within 15-32 degrees C at all times. Temperatures outside this range require immediate action (load reduction, emergency ventilation, or equipment shutdown if necessary to prevent damage).
- **R-D4-005-02:** Relative humidity must remain between 20-60% at all times. Persistent humidity outside this range accelerates equipment degradation.
- **R-D4-005-03:** Environmental monitoring must be continuous and logged. Gaps in monitoring data are gaps in the institution's knowledge of its equipment's operating conditions.
- **R-D4-005-04:** The environmental control audit must be performed at least semi-annually, timed to capture the most demanding seasonal conditions.
- **R-D4-005-05:** All filtration media must be inspected monthly and replaced on a schedule documented in the operational log. Clogged filters reduce cooling effectiveness and increase energy consumption.
- **R-D4-005-06:** The power consumption of all environmental control systems must be included in the load analysis for D4-002 (solar sizing) and D4-003 (battery sizing). Environmental controls that are not accounted for in the power budget will cause power shortfalls.
- **R-D4-005-07:** Equipment must not be powered on until it has acclimated to the room temperature. Cold equipment in a warm, humid room will form condensation. Allow a minimum of 2 hours for acclimation.

## 6. Failure Modes

- **Active cooling failure.** The air conditioner or mechanical cooling system stops operating. Impact: room temperature begins to rise. The rate of rise depends on the heat load and the room's thermal mass. A well-insulated, high-thermal-mass room may take hours to reach the allowable limit. A poorly insulated, low-thermal-mass room may reach it in 30-60 minutes. Detection: temperature sensor alerts. Mitigation: implement the cooling failure protocol (Section 7).

- **Dehumidifier failure.** In humid climates, the dehumidifier's failure allows relative humidity to climb. Impact: corrosion risk increases. Condensation risk increases if temperature drops. Detection: humidity sensor alerts. Mitigation: manual ventilation with dry outside air if available. Reduce moisture sources. Repair or replace the dehumidifier.

- **Ventilation blockage.** Intake or exhaust paths become blocked (filter clogging, debris, pest nesting, accidental obstruction). Impact: reduced airflow, rising temperature, potential positive/negative pressure loss. Detection: temperature rise without corresponding cooling system failure. Visual inspection. Mitigation: clear the blockage. Inspect all ventilation paths during the semi-annual audit.

- **Sensor failure.** A temperature or humidity sensor fails or becomes decalibrated. Impact: monitoring gap -- the institution loses awareness of conditions at that sensor's location. Detection: sudden flat-line or erratic readings in the sensor log. Comparison to adjacent sensors shows divergence. Mitigation: replace the sensor. Maintain spare sensors in inventory.

- **Water intrusion.** Leak from dehumidifier, air conditioner condensate, plumbing, or external source. Impact: potential for equipment damage if water contacts electronics. Electrical hazard. Humidity spike. Detection: water leak sensors, visual observation. Mitigation: immediate cleanup. Power down any equipment at risk of water contact. Identify and repair the water source.

- **Dust event.** Construction activity, natural disaster, or filter failure allows significant dust ingress. Impact: accelerated clogging of equipment fans and heat sinks, potential for electrical shorts from conductive dust. Detection: visible dust accumulation, temperature rise from clogged cooling. Mitigation: emergency cleaning. Identify and seal the ingress path. Replace compromised filters. Inspect and clean all affected equipment.

## 7. Recovery Procedures

1. **Cooling failure -- temperature rising:** If active cooling fails and passive ventilation cannot maintain temperature within the allowable range, implement load shedding. Shut down non-essential equipment in order of priority: development systems first, then non-critical services, then secondary storage, then primary services. The institution's shutdown priority list (documented in D4-006) determines the order. If temperature reaches 32 degrees C with all non-essential loads shed, shut down remaining equipment in an orderly fashion to prevent damage. Diagnose and repair the cooling system before resuming operations.

2. **Cooling failure -- extended duration:** If the cooling system cannot be repaired within 24 hours, implement temporary measures: portable fans to increase air circulation, temporary ventilation openings (filtered), reduction of computing load to the minimum required for data integrity. If necessary, operate equipment only during the coolest hours (overnight) and shut down during peak heat. Document the entire event.

3. **Humidity excursion -- high:** If humidity exceeds 60% and the dehumidifier is non-functional, increase ventilation with outside air if the outside air is drier. If outside air is also humid, reduce the ventilation rate to limit moisture ingress and wait for the dehumidifier repair. Monitor for condensation continuously. If condensation is observed on any equipment surface, shut that equipment down immediately.

4. **Water contact with equipment:** Power down the affected equipment immediately. Do not attempt to wipe water off energized electronics. Once power is removed, carefully dry all accessible surfaces. Allow the equipment to dry completely in a warm, well-ventilated environment for at least 48 hours. Before powering back on, inspect for corrosion or residue. Test the equipment thoroughly before returning it to production service.

5. **Post-dust-event recovery:** Shut down equipment in the affected area. Inspect all air intakes and internal components. Clean thoroughly with filtered compressed air. Replace any intake filters. Seal the dust ingress path. Monitor equipment temperatures closely for the following week to confirm that no hidden dust accumulation is causing elevated temperatures.

## 8. Evolution Path

- **Years 0-5:** Establish baseline environmental data. The first two years of continuous monitoring reveal the seasonal extremes and the effectiveness of the passive design. Adjust active systems based on actual data rather than design estimates.

- **Years 5-10:** Active cooling equipment (mini-splits, dehumidifiers) approaches mid-life. Refrigerant may need topping off. Filters may need replacement more frequently as they age. Source replacement units while current models are available. The passive design elements (insulation, thermal mass, ventilation paths) should remain stable and maintenance-free.

- **Years 10-15:** First generation of active cooling equipment may need replacement. This is an opportunity to upgrade to more efficient models. Evaluate the heat load: has it changed as computing equipment was replaced? Newer servers often consume less power and generate less heat than older ones. Resize the cooling system if appropriate.

- **Years 15-20:** Review the passive design. Has the building envelope changed? Has insulation degraded? Has vegetation growth affected ventilation paths? The semi-annual audit should catch these issues, but a comprehensive review at this interval is warranted.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Environmental control is where the off-grid constraint bites hardest. A grid-connected datacenter runs a 5,000-watt air conditioner without thinking about it. An off-grid institution must generate, store, and deliver every watt that the cooling system consumes. This is why the passive-first philosophy is not a nice-to-have; it is an economic necessity. Every degree of temperature rise that the building's passive design can handle without mechanical cooling is a degree's worth of air conditioning energy saved -- energy that can instead power the computing equipment that is the reason the institution exists.

If you are building the equipment room from scratch, invest in insulation and thermal mass. A below-grade room with concrete walls and good insulation may not need active cooling at all in temperate climates. If you are adapting an existing room, do what you can to improve insulation and ventilation before deploying an air conditioner. The air conditioner should be the last resort, not the first response.

One more thing: monitor your environment even when everything seems fine. Especially when everything seems fine. The temperature trend that moves from 22 degrees C to 24 degrees C over six months is telling you something. A clogged filter, a degrading refrigerant charge, an insulation failure, or a new heat source has entered the equation. Catch it at 24 degrees C and the fix is easy. Catch it at 35 degrees C and you are doing emergency shutdowns.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (off-grid mandate, self-built mandate)
- OPS-001 -- Operations Philosophy (maintenance tempo, documentation-first principle)
- D4-001 -- Infrastructure Architecture Overview
- D4-002 -- Solar Power System Design and Maintenance (power budget for cooling)
- D4-003 -- Battery Systems (power budget for cooling)
- D4-004 -- Network Hardware Lifecycle Management (environmental requirements for network equipment)
- D4-006 -- Power Distribution and UPS Management (power for environmental controls)
- ASHRAE TC 9.9 -- Thermal Guidelines for Data Processing Environments

---

---

# D4-006 -- Power Distribution and UPS Management

**Document ID:** D4-006
**Domain:** 4 -- Infrastructure & Power
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, OPS-001, D4-001, D4-002, D4-003
**Depended Upon By:** D4-004 (Network Hardware -- power supply), D4-005 (Environmental Controls -- power supply). All articles involving equipment power connections, power reliability, or shutdown procedures.

---

## 1. Purpose

This article provides a complete reference for the internal power distribution system of the holm.chat Documentation Institution: how electrical power flows from the battery bank and inverter to every piece of equipment, how uninterruptible power supply (UPS) systems provide bridging power during transitions and failures, how loads are balanced across circuits, and how the institution responds to extended power outages including generator integration.

The solar array generates power. The battery bank stores it. The inverter converts it to usable AC. But between the inverter's output terminals and the computing equipment's power supplies, there is an entire distribution system that must be designed, documented, and maintained. Circuit breakers, distribution panels, UPS units, power distribution units (PDUs), transfer switches, and the wiring that connects them all -- each component is a potential point of failure, and a failure in the power distribution system can damage equipment or cause data loss just as surely as a failure in the generation system.

This article is also the institution's guide to power emergencies. When the sun does not shine for a week, when the inverter fails, when the battery bank is depleted, or when a circuit breaker trips at 3 AM -- the procedures are here. The operator who faces a power emergency should not have to reason from first principles. They should have to follow documented steps.

## 2. Scope

This article covers:

- Power distribution architecture: from inverter output to equipment power supply input.
- Circuit design: how many circuits, what capacity, how they are protected.
- UPS systems: sizing, battery maintenance, testing, and replacement.
- Transfer switch configuration: automatic and manual transfer between power sources.
- Load balancing: distributing equipment across circuits to prevent overload.
- The power audit: a structured assessment of the entire power distribution chain.
- Power monitoring: what to measure, how to measure it, and how to interpret the data.
- Extended power outage procedures: load shedding, orderly shutdown, and generator integration.
- Generator systems: selection, fuel management, connection, and operational procedures.

This article does not cover solar panel systems (see D4-002) or battery bank management (see D4-003) except where they interface with the distribution system. It covers power from the point where it leaves the battery/inverter system to the point where it enters the equipment.

## 3. Background

### 3.1 The Power Distribution Chain

In a grid-connected building, power distribution is handled by the utility and the building's electrical panel. The operator plugs equipment in and gives it no further thought. In an off-grid institution, the operator is the utility. The entire power distribution chain is the operator's responsibility: from solar panel to battery to inverter to distribution panel to circuit breaker to outlet to power strip to equipment power supply. A failure anywhere in this chain interrupts power to everything downstream.

The chain has another property that demands attention: it concentrates risk. A single inverter failure takes out all AC loads. A single distribution panel fault takes out all circuits it feeds. A single circuit breaker trip takes out everything on that circuit. The distribution system must be designed with these concentration points in mind, and redundancy must be introduced where the consequences of failure are severe.

### 3.2 Why UPS in an Off-Grid System

The question is natural: if the entire institution runs on batteries already, why add UPS batteries? The answer is about transition time.

The institution's main battery bank powers the inverter, which produces AC power. But the inverter itself can fail, be switched off for maintenance, or trip on a fault condition. When it does, AC power disappears instantly. Computing equipment -- especially servers and storage arrays with data in flight -- can suffer data corruption from unclean shutdowns. A UPS provides bridging power: enough time for the operator to respond, switch to a backup inverter, start a generator, or perform an orderly shutdown of critical systems.

The UPS also provides power conditioning: it filters electrical noise, suppresses voltage spikes, and ensures that the power reaching the equipment is stable and clean. Even a pure sine wave inverter can produce transient anomalies during load changes or switching events. The UPS smooths these out.

### 3.3 The Generator Question

A generator is the power source of last resort. When the solar array cannot generate, the battery bank is depleted, and the institution faces a choice between shutting down and running a combustion engine -- the generator is the answer. Not every institution needs one. In climates with reliable solar resources and modest computing loads, a properly sized solar and battery system may provide sufficient autonomy year-round. In climates with extended low-sun periods (northern latitudes in winter), a generator provides the margin of safety that solar and batteries alone cannot.

This article covers generator integration as an optional but recommended component of the power system. The decision to include a generator is a Tier 2 governance decision based on the site's solar resource, the institution's load profile, and the operator's risk tolerance.

## 4. System Model

### 4.1 Power Distribution Architecture

The institution's power distribution follows a hierarchical architecture:

**Level 1: Source.** The primary source is the inverter, fed by the battery bank (which is charged by the solar array). The secondary source is the generator (if installed). A transfer switch -- manual or automatic -- selects between the primary and secondary sources.

**Level 2: Main distribution panel.** The transfer switch output feeds a main distribution panel (breaker panel). The panel contains circuit breakers sized for each branch circuit. All wiring from this point forward is standard AC electrical wiring installed to applicable electrical codes.

**Level 3: Branch circuits.** Individual circuits run from the distribution panel to the equipment areas. Each circuit serves a defined set of outlets or equipment. Circuits are segregated by function: critical computing equipment on dedicated circuits, environmental controls on separate circuits, lighting and general power on separate circuits.

**Level 4: UPS systems.** Critical computing equipment is powered through UPS units. The UPS is plugged into the branch circuit and the equipment is plugged into the UPS. The UPS provides bridging power and power conditioning between the branch circuit and the equipment.

**Level 5: Power distribution units (PDUs).** In rack-mounted installations, PDUs distribute power from the UPS to individual devices in the rack. PDUs with per-outlet monitoring allow the operator to measure the power consumption of individual devices.

### 4.2 Circuit Design

**Sizing:** Each branch circuit must be rated for 125% of the maximum expected continuous load on that circuit (per standard electrical code practice). A circuit serving 800 watts of continuous load requires a circuit breaker rated for at least 1,000 watts (approximately 8.3 amps at 120VAC or 4.3 amps at 230VAC).

**Segregation:** Separate critical and non-critical loads onto different circuits. This ensures that a non-critical fault (a tripped breaker on the lighting circuit, for example) does not affect computing equipment. Recommended circuit groups:

- Critical computing (servers, primary storage): dedicated circuits through UPS.
- Network infrastructure (switches, access points): dedicated circuit through UPS.
- Secondary storage and backup systems: dedicated circuit, optionally through UPS.
- Environmental controls (cooling, dehumidifier): dedicated circuit(s), not through UPS (too high a load for battery-backed UPS, and brief power interruptions do not damage cooling equipment).
- Lighting and general power: separate circuit, no UPS.
- Workstation: separate circuit, optionally through UPS.

**Labeling:** Every circuit breaker must be labeled with the circuit number and the equipment it serves. Every outlet must be labeled with the circuit number. The circuit map -- showing which circuit serves which outlet and which equipment -- must be documented and posted in the electrical room.

### 4.3 UPS Sizing and Selection

**Capacity sizing:** The UPS must support the connected load for the desired runtime. Determine the total load in watts (measured, not estimated). Multiply by the desired runtime in hours. Add a 25% margin. This gives the required UPS capacity in watt-hours. Example: 500 watts of critical computing equipment, 15 minutes (0.25 hours) desired runtime: 500 x 0.25 x 1.25 = 156 Wh minimum UPS capacity.

Note that 15 minutes is a minimum bridging time for the operator to assess the situation and initiate an orderly shutdown. If the operator may not be present when a power event occurs (overnight, for example), longer runtimes or automated shutdown scripts are necessary.

**Type selection:** Online (double-conversion) UPS is preferred for critical computing loads. It continuously converts AC to DC and back to AC, providing complete isolation from input power anomalies. Line-interactive UPS is acceptable for less critical loads. Standby (offline) UPS is not recommended for computing equipment due to the transfer time gap when switching to battery.

**Battery type:** UPS units typically use sealed lead-acid (SLA/VRLA) batteries internally. These batteries have a shorter life than the main battery bank (3-5 years typical) and must be replaced on schedule. Some higher-end UPS units accept external battery packs for extended runtime. LFP-based UPS systems are becoming available and offer longer battery life; evaluate if available.

**Output waveform:** Pure sine wave output only. This requirement is the same as for the main inverter (D4-002) and for the same reasons.

### 4.4 UPS Battery Maintenance

UPS batteries are the weakest link in the power distribution chain. They are small, sealed lead-acid batteries that degrade faster than the main battery bank because they are often subjected to higher temperatures (inside the UPS enclosure) and receive less attention.

**Monthly:** Check UPS status indicators. Record battery charge level and estimated runtime. Compare estimated runtime to the rated runtime -- declining runtime indicates battery degradation.

**Quarterly:** Run a UPS self-test (most UPS units have a built-in self-test function). Record the test result. If the UPS reports a battery warning, schedule replacement.

**Annually:** Perform a load test: run the UPS on battery for a measured duration and compare to the expected runtime. If actual runtime is less than 80% of rated runtime with fresh batteries, replace the batteries.

**Replacement schedule:** Replace UPS batteries every 3-4 years regardless of test results. Lead-acid batteries in UPS service have a predictable lifespan, and the cost of proactive replacement is trivial compared to the cost of a UPS battery failure during a power event. Order replacement batteries before removing old ones. Test replacement batteries immediately after installation.

### 4.5 Transfer Switch

The transfer switch selects between the primary power source (inverter) and the secondary power source (generator). It can be manual or automatic.

**Manual transfer switch:** The operator physically switches between sources. Simple, reliable, no electronics to fail. Disadvantage: requires the operator to be present and aware that a transfer is needed.

**Automatic transfer switch (ATS):** Monitors the primary source and automatically switches to the secondary source when the primary fails. Returns to the primary source when it is restored. Advantage: responds faster than a human. Disadvantage: adds complexity, requires power, is itself a potential failure point.

For the institution, an ATS is recommended if a generator is installed. The transfer sequence: primary source fails, UPS absorbs the load, ATS starts the generator, ATS verifies generator output is stable, ATS transfers the load to the generator, UPS returns to standby. When the primary source is restored, the ATS transfers back.

**Transfer switch testing:** Test the transfer switch quarterly by simulating a primary source failure. Verify that the transfer occurs cleanly with no disruption to UPS-backed loads. Document the test results.

### 4.6 Load Balancing

In a multi-circuit installation, distribute loads across circuits to prevent any single circuit from being overloaded while others are underutilized. The goal is to keep each circuit below 80% of its rated capacity during normal operation, leaving headroom for transient peaks and future load additions.

**The power audit:** The power audit is a comprehensive measurement of the entire power distribution chain, performed annually. It measures:

- Total power consumption at the inverter output (total institutional load).
- Power consumption on each branch circuit.
- Power consumption of each major device (using a clamp meter or per-outlet PDU monitoring).
- UPS input power, output power, and efficiency (to verify UPS health).
- Voltage at multiple points in the distribution chain (to detect voltage drop from undersized wiring or poor connections).
- Temperature of all distribution panel breakers, wire terminations, and UPS units (hot spots indicate high-resistance connections or overloaded circuits).

Document the power audit results. Compare to previous audits. Track trends: is total consumption growing? Is any circuit approaching its capacity? Is the UPS runtime declining?

### 4.7 Extended Power Outage Procedures

An extended power outage occurs when the battery bank is depleted and no solar generation is available. The institution's response follows a staged protocol:

**Stage 1: Awareness (battery bank above 50% state of charge).** Normal operation. Monitor the weather forecast and solar generation trend. If generation is expected to be low for multiple days, begin voluntary load reduction: defer non-essential computing tasks, shut down development environments, reduce lighting.

**Stage 2: Conservation (battery bank at 30-50% state of charge).** Mandatory load shedding. Shut down all non-essential systems. Operate only critical systems: primary storage, network core, one workstation. Reduce environmental controls to minimum (rely on passive cooling and ventilation). Start the generator if available and if the outage is expected to continue.

**Stage 3: Critical (battery bank at 15-30% state of charge).** Operate only the systems required to maintain data integrity. Primary storage must remain online. Network can be shut down. Workstation can be shut down. If the generator is running, it should be carrying the full load and the battery bank should be recovering. If no generator is available, prepare for orderly shutdown.

**Stage 4: Emergency shutdown (battery bank at 10-15% state of charge).** Perform an orderly shutdown of all remaining systems. Ensure all data is flushed to disk, all filesystems are cleanly unmounted, and all systems are powered off gracefully. The battery bank's remaining capacity is reserved for the UPS to prevent abrupt shutdown of the last systems.

**Stage 5: Dark institution.** All systems are off. The battery bank is depleted. The institution waits for solar generation to resume or for the generator to be started. No data should be at risk because Stage 4 performed an orderly shutdown. When power is restored, bring systems up in reverse priority order: storage first, then network, then servers, then workstations, then non-essential systems.

Document the shutdown state of each system during the outage. When power is restored, verify data integrity before resuming normal operations.

### 4.8 Generator Integration

**Generator selection:** For institutional use, a diesel or propane generator in the 3-10 kW range is typical, depending on the institution's load profile. Sizing: the generator must support the institution's critical load plus the battery charging load plus environmental control load simultaneously. Oversize by 25% to avoid running the generator at continuous full load, which accelerates wear.

**Fuel management:** Store fuel in approved containers in a well-ventilated, fire-safe location away from the equipment room. Diesel fuel has a shelf life of 6-12 months without stabilizer, 12-24 months with stabilizer. Propane has an indefinite shelf life, making it advantageous for standby generator fuel. Maintain a fuel inventory log. Rotate fuel stock.

**Connection:** The generator connects to the power system through the transfer switch (Section 4.5). Never connect a generator directly to the institution's wiring without a transfer switch. Direct connection creates the risk of backfeed, which is both a safety hazard and an equipment damage risk.

**Starting procedure:** Start the generator according to the manufacturer's procedure. Allow the generator to warm up and stabilize (typically 2-5 minutes). Verify voltage and frequency output with a meter. Then engage the transfer switch to transfer the load.

**Runtime management:** Generators require periodic rest and maintenance. For extended outages, establish a run schedule: run the generator for 8-12 hours to charge batteries and run equipment, then shut down for 4-8 hours while batteries carry the load. This extends fuel supply and reduces generator wear.

**Maintenance:** Generators require regular maintenance even when not in use: monthly engine run (15-30 minutes under load), annual oil change, annual fuel filter replacement, spark plug or glow plug inspection, and air filter inspection. A generator that has sat idle for a year without maintenance may not start when needed. Treat generator maintenance as part of the monthly operations cycle (OPS-001).

## 5. Rules & Constraints

- **R-D4-006-01:** Every circuit breaker must be labeled, and the circuit map must be documented and kept current. An unlabeled breaker panel is a diagnostic nightmare during a power emergency.
- **R-D4-006-02:** Critical computing equipment must be powered through a UPS. Direct connection of critical equipment to branch circuits without UPS protection is prohibited.
- **R-D4-006-03:** UPS batteries must be replaced on a 3-4 year schedule regardless of test results. Aging UPS batteries are the most common cause of UPS failure during power events.
- **R-D4-006-04:** The power audit must be performed annually. Results must be documented and trends tracked.
- **R-D4-006-05:** The transfer switch must be tested quarterly by simulating a primary source failure. Untested transfer switches have an unquantified failure rate.
- **R-D4-006-06:** No circuit shall be loaded above 80% of its rated capacity during normal operation. Circuits approaching 80% must be rebalanced.
- **R-D4-006-07:** If a generator is installed, it must be exercised monthly under load. A generator that has not been run in 30 days has unknown readiness.
- **R-D4-006-08:** The extended power outage protocol (Section 4.7) must be printed on physical media and posted in the equipment room. During a power outage, the electronic documentation system may not be accessible.

## 6. Failure Modes

- **Inverter failure.** The main inverter stops producing AC power. Impact: total loss of AC power to the distribution system. UPS units absorb their connected loads; non-UPS loads go dark immediately. Detection: immediate -- non-UPS loads fail, UPS units alarm. Mitigation: switch to generator (via transfer switch) or to spare inverter (D4-002). UPS runtime provides the bridging time.

- **UPS battery failure.** The UPS batteries are depleted or have failed. The UPS cannot bridge a power event. Impact: when the next power interruption occurs, critical equipment experiences an unclean shutdown. Detection: UPS self-test failure, declining runtime measurements in quarterly tests. Mitigation: replace batteries on schedule. Do not defer replacement.

- **Circuit breaker trip.** A breaker trips due to overcurrent (overloaded circuit) or ground fault. Impact: all equipment on that circuit loses power. Detection: immediate for monitored circuits. May go undetected for unmonitored circuits until the operator notices a system is down. Mitigation: do not overload circuits. Investigate and resolve the cause before resetting the breaker. A breaker that trips repeatedly indicates a wiring fault, an equipment fault, or an overloaded circuit -- not a defective breaker.

- **Transfer switch failure.** The transfer switch fails to transfer to the secondary source during a primary source failure. Impact: the generator is running and ready, but the load remains on the failing primary source. Detection: generator is running but loads remain off or on UPS battery. Mitigation: manual bypass -- most transfer switches include a manual override. Test the transfer switch quarterly to detect failures before they matter.

- **Generator failure to start.** The generator does not start when needed. Impact: no secondary power source during an extended primary source outage. The institution relies entirely on battery reserves. Detection: immediate -- the generator does not respond to start command. Mitigation: monthly exercise under load. Maintain fuel freshness. Keep a documented troubleshooting guide for the specific generator model.

- **Wiring failure.** A connection in the distribution system fails due to corrosion, vibration, or thermal cycling. Impact: intermittent or complete loss of power on the affected circuit. Can cause arcing, which is a fire hazard. Detection: intermittent equipment reboots on a specific circuit, hot spots detected during power audit, tripped breaker with no apparent overload. Mitigation: annual thermal scan of all terminations during the power audit. Retorque all connections annually.

## 7. Recovery Procedures

1. **Inverter failure with generator available:** UPS units bridge the load. Start the generator. Verify generator output. Engage the transfer switch. Verify loads have transferred. Diagnose and repair or replace the inverter. When the inverter is repaired, test it off-load, then transfer back from the generator.

2. **Inverter failure without generator:** UPS units bridge the load for their rated runtime. Immediately begin orderly shutdown of non-critical systems. If a spare inverter is available (D4-002), install it. If no spare is available, perform a complete orderly shutdown before UPS batteries are depleted. The institution remains dark until the inverter is repaired or replaced.

3. **UPS failure during power event:** Equipment on the failed UPS experiences unclean shutdown. After primary power is restored, check all affected systems for data integrity (filesystem checks, database integrity verification, log analysis). Replace the UPS. Restore affected systems from backup if integrity checks fail.

4. **Circuit breaker trip -- cause known:** If the cause is a known temporary overload (equipment startup surge), reset the breaker and monitor. If the cause is a permanent overload, redistribute loads before resetting. If the cause is a ground fault, do not reset until the fault is identified and repaired.

5. **Extended outage recovery (Stage 5 to normal operation):** When power is restored (sun returns, generator running, inverter repaired), do not bring everything up at once. Start with storage systems. Verify storage integrity. Start network infrastructure. Verify network connectivity. Start servers one at a time. Verify each server's services. Start workstations. Resume non-essential systems. Resume environmental controls. This staged approach prevents a startup surge from tripping breakers or overloading the recovering power system, and it allows integrity verification at each step.

6. **Generator failure during extended outage:** If the generator fails during an extended outage and cannot be immediately repaired, implement the battery conservation protocol (Section 4.7, Stages 2-4). If battery reserves are already low, proceed directly to orderly shutdown. The generator should have a documented troubleshooting guide for common failure modes (fuel starvation, dead starter battery, clogged fuel filter, loss of spark/compression).

## 8. Evolution Path

- **Years 0-5:** The power distribution system is new. Establish baseline measurements with the first power audit. The initial load balance may need adjustment as the institution's equipment complement stabilizes. UPS batteries are new and reliable. Generator (if installed) requires only routine maintenance.

- **Years 5-10:** First UPS battery replacement cycle occurs. The distribution panel and wiring remain stable. Generator may need its first significant maintenance (carburetor cleaning for gasoline, injector service for diesel). Evaluate whether the circuit layout still matches the institution's needs as equipment has been added or moved.

- **Years 10-15:** Second or third UPS battery replacement. Consider replacing the UPS units themselves if they show signs of age (capacitor bulging, fan noise, erratic behavior). The distribution wiring should be inspected for insulation degradation, particularly in areas with temperature extremes. Generator may need overhaul or replacement depending on total runtime hours.

- **Years 15-20:** The distribution system is mature. The transfer switch, if mechanical, may need servicing (contact cleaning, spring replacement). Circuit breakers have a long life but should be exercised (tripped and reset) annually to prevent mechanical seizure. Evaluate whether the overall power architecture still serves the institution's needs or whether a redesign is warranted.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Power distribution is the most mundane and the most consequential infrastructure topic. Nobody is excited about circuit breakers. Nobody wants to spend a Saturday testing the transfer switch. But the institution's relationship with its power distribution system is like a building's relationship with its foundation: you do not think about it when it is working, and you cannot think about anything else when it is not.

The single most important practice in this article is the power audit. Measure everything. Know your loads. Know your margins. Know which circuits are approaching capacity and which have headroom. A power audit takes a few hours once a year. A power fire takes your institution away permanently.

The second most important practice is UPS battery replacement. UPS batteries fail silently. The UPS sits there with its green light on, looking perfectly healthy, while its batteries quietly turn into expensive paperweights. When the next power event arrives, the UPS transfers to battery, the battery delivers nothing, and your servers crash. Test your UPS. Replace the batteries on schedule. This is not optional.

On generators: if your climate has seasons where solar generation drops significantly, a generator is not optional either. Three days of battery autonomy sounds like plenty until you experience two weeks of overcast winter weather. A small propane generator, exercised monthly and maintained properly, is cheap insurance against the scenario where the sun simply does not cooperate. Propane does not go stale. A propane generator that has been sitting for six months will start just as readily as one that ran last week, provided it has been properly maintained. That property alone makes propane the preferred generator fuel for standby applications.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (off-grid mandate)
- OPS-001 -- Operations Philosophy (maintenance tempo, power audit as part of annual review)
- D4-001 -- Infrastructure Architecture Overview
- D4-002 -- Solar Power System Design and Maintenance (primary power source)
- D4-003 -- Battery Systems: Selection, Management, and End-of-Life (energy storage)
- D4-004 -- Network Hardware Lifecycle Management (power requirements for network equipment)
- D4-005 -- Environmental Control Systems (power requirements for environmental controls)
- NFPA 70 (NEC) -- National Electrical Code (applicable wiring standards)
- IEEE 1100 -- Recommended Practice for Powering and Grounding Sensitive Electronic Equipment
- NFPA 110 -- Standard for Emergency and Standby Power Systems

---

*End of Stage 4: Specialized Systems -- Infrastructure Advanced*

**Document Total:** 5 articles
**Combined Estimated Word Count:** ~15,000 words
**Status:** All five articles ratified as of 2026-02-16.
**Next Stage:** Stage 4 continues with additional specialized system articles for storage, software platforms, and data management.
