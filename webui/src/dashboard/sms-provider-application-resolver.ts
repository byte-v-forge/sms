import { listSmsApplications } from './sms-api';
import { applicationChoices, matchApplicationChoice, providerApplicationKey, type ApplicationChoice } from './sms-compare-options';

export async function resolveSmsProviderApplicationKey(providerKey: string, serviceText: string, application: ApplicationChoice | undefined) {
  const selectedKey = application?.providerApplicationKeys[providerKey];
  if (selectedKey) return selectedKey;
  const response = await listSmsApplications({ providerKeys: [providerKey], searchText: serviceText });
  const choices = applicationChoices([{ providerKey, applications: response.applications || [], error: response.error }], []);
  const choice = matchApplicationChoice(serviceText, choices) || singleChoice(choices);
  return providerApplicationKey(choice, providerKey, serviceText);
}

function singleChoice(choices: ApplicationChoice[]) {
  return choices.length === 1 ? choices[0] : undefined;
}
