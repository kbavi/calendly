# Calendly

**Product Requirements Document (PRD)**
---

### **1. Summary**

Calendly aims to define simplicity and efficiency in scheduling. It provides an intuitive platform that integrates calendars, event types, events, availabilities, and overlap discovery. This MVP focuses on delivering core functionality to help users schedule and reschedule meetings with minimal friction.

---

### **2. Objectives**

- **Core Goals:**
  - Enable users to connect and manage their calendar.
  - Allow users to define and share their availability.
  - Enable finding availability overlaps between participants.
  - Allow booking meetings.
  - Allow effortless rescheduling of existing meetings.
---

### **3. Key Features and Requirements**

#### **3.1 Calendars**
- **Description:** Users can connect and view their calendar.
- **Key Functionalities:**
  - Support for a single calendar integration (Mock Google Calendar for now). Future versions will expand this to connect multiple calendars across providers.

#### **3.3 Events**
- **Description:** Management of all scheduled events.
- **Key Functionalities:**
  - Event details (participants, notes).
  - Allow quick booking of 30, 60, 90 minute events.

#### **3.4 Availabilities**
- **Description:** Users can set their availability for event bookings.
- **Key Functionalities:**
  - Configurable time blocks.
  - Allow users to set date specific availability (e.g. 10am to 12pm on 1st Dec 2024).
  - Allow users to set day-of-the-week-specific availability (e.g. 10am to 12pm on every Monday).
  - Allow multiple availabilities for a day (e.g. 10am to 12pm and 2pm to 5pm, Monday to Friday).

#### **3.5 Scheduling**
- **Description:** Enables seamless discovery of mutual availability between participants.
- **Key Functionalities:**
  - Algorithms to find the earliest mutually available time slot.
  - Efficiently find overlaps between availabilities of multiple participants between a time range.
  - One-click scheduling for invitees via booking links.
  - Booked or cancelled events should be reflected in the availability of the participants.

---

### **4. User Stories**

#### **4.1 As a user,**
- I want to connect my Calendar to Calendly so that my events sync in real-time.
- I want to create a 30-minute consultation event type so I can schedule calls with clients.
- I want to define my availability so that I’m not booked during personal time.
- I want to share my scheduling link to let others book a time with me effortlessly.
- I want to reschedule an event easily if my plans change.

---

### **5. Assumptions and Constraints**

- No support for enterprise-level features like SSO in the MVP.
- No support for more than one calendar in the MVP.
- Platform supports only English language initially.
- MVP does not support any notifications but the platform should be able to easily integrate with notifications in future.
- The MVP does not support Recurring Events
- The MVP does not support multiple timezones. All events should be considered to be in UTC time.
- User Authentication is not supported in the MVP. The responsibility of privacy and security is left to the caller.

### **6. Test Cases**

#### **6.1 User Availability Patterns**

| User Type | Availability Pattern | Description |
|-----------|---------------------|-------------|
| Weekday User | Monday-Friday | Available on weekdays only |
| Midweek User | Wednesday-Saturday | Available mid-week through weekend |
| Weekend User | Saturday-Sunday | Available on weekends only |

#### **6.2 Overlap Scenarios**

| Test Case | Participants | Time Period | Expected Behavior | Rationale |
|-----------|-------------|-------------|-------------------|-----------|
| All Users No Overlap | Weekday, Midweek, Weekend users | Monday 9 AM - 5 PM | No overlapping slots | No time period exists where all three users are available |
| Weekday-Midweek Friday Overlap | Weekday and Midweek users | Friday 9 AM - 5 PM | 2 overlapping slots | Both users are available on Fridays |
| Midweek-Weekend Saturday Overlap | Midweek and Weekend users | Saturday 9 AM - 5 PM | 1 overlapping slot | Both users are available on Saturdays |
| No Common Availability | Weekday and Weekend users | Saturday 9 AM - 5 PM | No overlapping slots | Users have mutually exclusive schedules |

#### **6.3 System Behavior Requirements**

1. **Availability Rules**
   - System must correctly handle multiple availability rules per user
   - Each user type should have the correct number of availability rules:
     - Weekday user: 5 rules (Mon-Fri)
     - Midweek user: 4 rules (Wed-Sat)
     - Weekend user: 2 rules (Sat-Sun)

2. **Overlap Detection**
   - System must efficiently find overlapping availabilities between multiple users
   - Must consider both:
     - User-defined availability rules
     - Existing calendar events

3. **Time Handling**
   - All times must be handled in UTC
   - System must properly handle day-of-week based rules
   - Must support time ranges within a single day (e.g., 9 AM - 5 PM)


#### **6.4 Additional Test Scenarios**

##### **A. Working Hours Patterns**

| User Type | Availability Pattern | Description |
|-----------|---------------------|-------------|
| Early Bird | Monday-Friday, 6 AM - 2 PM | Early shift worker |
| Night Owl | Monday-Friday, 2 PM - 10 PM | Late shift worker |
| Split Shift | Monday-Friday, 9 AM - 1 PM, 4 PM - 8 PM | Break in middle of day |
| Flexible Hours | Monday-Friday, varying blocks | Different hours each day |

**Test Cases:**
1. Early Bird + Regular (9-5) Overlap
2. Night Owl + Regular (9-5) Overlap
3. Split Shift + Regular (9-5) Partial Overlap
4. Multiple Split Shift Users Overlap

##### **B. Partial Availability Patterns**

| User Type | Availability Pattern | Description |
|-----------|---------------------|-------------|
| Part-Time | Mon, Wed, Fri, 9 AM - 3 PM | Limited weekdays |
| Morning Only | Monday-Friday, 7 AM - 12 PM | Morning availability |
| Afternoon Only | Monday-Friday, 1 PM - 6 PM | Afternoon availability |
| Custom Days | Tue, Thu, Sat, 10 AM - 6 PM | Specific days only |

**Test Cases:**
1. Part-Time + Full-Time Overlap
2. Morning Only + Afternoon Only No Overlap
3. Custom Days + Weekend User Overlap
4. Multiple Part-Time Users Common Slots

##### **C. Edge Cases**

| Scenario | Description | Expected Behavior |
|----------|-------------|------------------|
| Boundary Time | Meeting at availability boundary (e.g., 5 PM) | Should consider end time in availability check |
| Minimum Duration | 15-min slots between different availability patterns | Should find minimal overlapping windows |
| Maximum Participants | 10+ users with different patterns | Should efficiently find common slots |
| Multi-Day Event | Event spanning multiple days | Should handle events spanning multiple days |

**Test Cases:**
1. Boundary Time Overlap Detection
2. Short Duration Multi-User Overlap
3. Large Group Common Availability
4. Holiday Period Availability

##### **D. System Load Patterns**

| Test Type | Description | Expected Performance |
|-----------|-------------|---------------------|
| Large Date Range | 6-month availability window | Efficient processing |
| Multiple Rules | User with 20+ availability rules | Quick rule evaluation |
| Dense Calendar | Calendar with 100+ existing events | Fast overlap detection |

**Test Cases:**
1. Concurrent User Requests
2. Extended Time Range Search
3. Complex Rule Combination
4. High Event Density Processing

These additional test suites ensure:
- Support for diverse working patterns
- Efficient handling of time zone differences
- Proper management of partial availability
- System performance under various loads
- Correct handling of edge cases and boundary conditions

