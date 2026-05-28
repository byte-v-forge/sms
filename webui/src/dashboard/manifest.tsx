import { MessageSquareText } from 'lucide-react';
import { DashboardNavSection, type DashboardModuleRegistration } from '@byte-v-forge/common-ui';
import { SmsPage } from './sms-page';

const registration: DashboardModuleRegistration = {
  manifest: {
    id: 'sms',
    nav: [
      {
        key: 'sms',
        label: 'SMS',
        icon: 'sms',
        section: DashboardNavSection.DASHBOARD_NAV_SECTION_INFRASTRUCTURE,
        required_services: ['sms-service'],
        order: 30
      }
    ]
  },
  icons: {
    sms: <MessageSquareText size={17} />
  },
  views: {
    sms: () => <SmsPage />
  }
};

export default registration;
