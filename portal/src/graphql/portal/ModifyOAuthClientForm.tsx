import React from "react";
import cn from "classnames";
import produce from "immer";
import { DirectionalHint } from "@fluentui/react";

import LabelWithTooltip from "../../LabelWithTooltip";
import FormTextField from "../../FormTextField";
import {
  useIntegerTextField,
  useTagPickerWithNewTags,
  useTextField,
} from "../../hook/useInput";
import { OAuthClientConfig } from "../../types";
import { ensureNonEmptyString } from "../../util/misc";

import styles from "./ModifyOAuthClientForm.module.scss";
import FormTagPicker from "../../FormTagPicker";

interface ModifyOAuthClientFormProps {
  className?: string;
  clientConfig: OAuthClientConfig;
  onClientConfigChange: (newClientConfig: OAuthClientConfig) => void;
}

const jsonPointerRegExp: Pick<
  Record<keyof OAuthClientConfig, { field: RegExp; parent: RegExp }>,
  | "name"
  | "access_token_lifetime_seconds"
  | "refresh_token_lifetime_seconds"
  | "redirect_uris"
  | "post_logout_redirect_uris"
> = {
  name: {
    field: /^\/oauth\/clients\/[0-9]+\/name$/,
    parent: /^\/oauth\/clients\/[0-9]+$/,
  },
  access_token_lifetime_seconds: {
    field: /^\/oauth\/clients\/[0-9]+\/access_token_lifetime_seconds$/,
    parent: /^\/oauth\/clients\/[0-9]+$/,
  },
  refresh_token_lifetime_seconds: {
    field: /^\/oauth\/clients\/[0-9]+\/refresh_token_lifetime_seconds$/,
    parent: /^\/oauth\/clients\/[0-9]+$/,
  },
  redirect_uris: {
    field: /^\/oauth\/clients\/[0-9]+\/redirect_uris(\/[0-9]+)?$/,
    parent: /^\/oauth\/clients\/[0-9]+/,
  },
  post_logout_redirect_uris: {
    field: /^\/oauth\/clients\/[0-9]+\/post_logout_redirect_uris(\/[0-9]+)?$/,
    parent: /^\/oauth\/clients\/[0-9]+/,
  },
};

export function getReducedClientConfig(
  clientConfig: OAuthClientConfig
): Omit<OAuthClientConfig, "grant_types" | "response_types"> {
  const {
    grant_types: _grant_types,
    response_types: _response_types,
    ...rest
  } = clientConfig;

  return {
    ...rest,
    post_logout_redirect_uris: rest.post_logout_redirect_uris ?? [],
  };
}

function updateClientConfig<K extends keyof OAuthClientConfig>(
  clientConfig: OAuthClientConfig,
  field: K,
  newValue: OAuthClientConfig[K]
): OAuthClientConfig {
  return produce(clientConfig, (draftConfig) => {
    draftConfig[field] = newValue;
  });
}

function convertIntegerStringToNumber(value: string): number | undefined {
  // Number("") = 0
  const numericValue = Number(value);
  return numericValue === 0 ? undefined : numericValue;
}

const ModifyOAuthClientForm: React.FC<ModifyOAuthClientFormProps> = function ModifyOAuthClientForm(
  props: ModifyOAuthClientFormProps
) {
  const { className, clientConfig, onClientConfigChange } = props;

  const { onChange: onClientNameChange } = useTextField((value) => {
    onClientConfigChange(
      updateClientConfig(clientConfig, "name", ensureNonEmptyString(value))
    );
  });

  const { onChange: onAccessTokenLifetimeChange } = useIntegerTextField(
    (value) => {
      onClientConfigChange(
        updateClientConfig(
          clientConfig,
          "access_token_lifetime_seconds",
          convertIntegerStringToNumber(value)
        )
      );
    }
  );
  const { onChange: onRefreshTokenLifetimeChange } = useIntegerTextField(
    (value) => {
      onClientConfigChange(
        updateClientConfig(
          clientConfig,
          "refresh_token_lifetime_seconds",
          convertIntegerStringToNumber(value)
        )
      );
    }
  );

  const {
    selectedItems: redirectUriItems,
    onChange: onRedirectUrisChange,
    onResolveSuggestions: onResolveRedirectUrisSuggestions,
  } = useTagPickerWithNewTags(clientConfig.redirect_uris, (list) => {
    onClientConfigChange(
      updateClientConfig(clientConfig, "redirect_uris", list)
    );
  });

  const {
    selectedItems: postLogoutRedirectUriItems,
    onChange: onPostLogoutRedirectUrisChange,
    onResolveSuggestions: onResolvePostLogoutUrisSuggestions,
  } = useTagPickerWithNewTags(
    clientConfig.post_logout_redirect_uris ?? [],
    (list) => {
      onClientConfigChange(
        updateClientConfig(
          clientConfig,
          "post_logout_redirect_uris",
          list.length > 0 ? list : undefined
        )
      );
    }
  );

  return (
    <section className={cn(styles.root, className)}>
      <FormTextField
        jsonPointer={jsonPointerRegExp.name.field}
        parentJSONPointer={jsonPointerRegExp.name.parent}
        fieldName="name"
        fieldNameMessageID="ModifyOAuthClientForm.name-label"
        className={styles.inputField}
        value={clientConfig.name ?? ""}
        onChange={onClientNameChange}
        required={true}
      />
      <FormTextField
        jsonPointer={jsonPointerRegExp.access_token_lifetime_seconds.field}
        parentJSONPointer={
          jsonPointerRegExp.access_token_lifetime_seconds.parent
        }
        fieldName="access_token_lifetime_seconds"
        fieldNameMessageID="ModifyOAuthClientForm.acces-token-lifetime-label"
        className={styles.inputField}
        value={clientConfig.access_token_lifetime_seconds?.toString() ?? ""}
        onChange={onAccessTokenLifetimeChange}
      />
      <FormTextField
        jsonPointer={jsonPointerRegExp.refresh_token_lifetime_seconds.field}
        parentJSONPointer={
          jsonPointerRegExp.refresh_token_lifetime_seconds.parent
        }
        fieldName="refresh_token_lifetime_seconds"
        fieldNameMessageID="ModifyOAuthClientForm.refresh-token-lifetime-label"
        className={styles.inputField}
        value={clientConfig.refresh_token_lifetime_seconds?.toString() ?? ""}
        onChange={onRefreshTokenLifetimeChange}
      />
      <FormTagPicker
        jsonPointer={jsonPointerRegExp.redirect_uris.field}
        parentJSONPointer={jsonPointerRegExp.redirect_uris.parent}
        fieldName="redirect_uris"
        fieldNameMessageID="ModifyOAuthClientForm.redirect-uris-label"
        label={
          <LabelWithTooltip
            labelId="ModifyOAuthClientForm.redirect-uris-label"
            tooltipHeaderId="ModifyOAuthClientForm.redirect-uris-label"
            tooltipMessageId="ModifyOAuthClientForm.redirect-uris-tooltip-message"
            directionalHint={DirectionalHint.bottomLeftEdge}
            required={true}
          />
        }
        className={styles.inputField}
        onChange={onRedirectUrisChange}
        onResolveSuggestions={onResolveRedirectUrisSuggestions}
        selectedItems={redirectUriItems}
      />
      <FormTagPicker
        jsonPointer={jsonPointerRegExp.post_logout_redirect_uris.field}
        parentJSONPointer={jsonPointerRegExp.post_logout_redirect_uris.parent}
        fieldName="post_logout_redirect_uris"
        fieldNameMessageID="ModifyOAuthClientForm.post-logout-redirect-uris-label"
        label={
          <LabelWithTooltip
            labelId="ModifyOAuthClientForm.post-logout-redirect-uris-label"
            tooltipHeaderId="ModifyOAuthClientForm.post-logout-redirect-uris-label"
            tooltipMessageId="ModifyOAuthClientForm.post-logout-redirect-uris-tooltip-message"
            directionalHint={DirectionalHint.bottomLeftEdge}
          />
        }
        className={styles.inputField}
        onChange={onPostLogoutRedirectUrisChange}
        onResolveSuggestions={onResolvePostLogoutUrisSuggestions}
        selectedItems={postLogoutRedirectUriItems}
      />
    </section>
  );
};

export default ModifyOAuthClientForm;
