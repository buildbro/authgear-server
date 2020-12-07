import React, { useCallback, useContext, useMemo, useState } from "react";
import {
  DefaultButton,
  Dialog,
  DialogFooter,
  ICommandBarItemProps,
  MessageBar,
  MessageBarType,
  PrimaryButton,
} from "@fluentui/react";
import { Context, FormattedMessage } from "@oursky/react-messageformat";
import { useValidationError } from "./error/useValidationError";
import { FormContext } from "./error/FormContext";
import ShowUnhandledValidationErrorCause from "./error/ShowUnhandledValidationErrorCauses";
import { useSystemConfig } from "./context/SystemConfigContext";
import ShowError from "./ShowError";
import NavigationBlockerDialog from "./NavigationBlockerDialog";
import CommandBarContainer from "./CommandBarContainer";
import { GenericErrorHandlingRule, parseError } from "./error/useGenericError";
import { APIError } from "./error/error";

export interface FormModel {
  updateError: unknown;
  isDirty: boolean;
  isUpdating: boolean;
  reset: () => void;
  save: () => void;
}

export interface SaveButtonProps {
  labelId: string;
  iconName: string;
}

export interface FormContainerProps {
  form: FormModel;
  canSave?: boolean;
  saveButtonProps?: SaveButtonProps;
  localError?: APIError | null;
  errorParseRules?: GenericErrorHandlingRule[];
  fallbackErrorMessageId?: string;
  farItems?: ICommandBarItemProps[];
  messageBar?: React.ReactNode;
}

const FormContainer: React.FC<FormContainerProps> = function FormContainer(
  props
) {
  const { updateError, isDirty, isUpdating, reset, save } = props.form;
  const {
    canSave = true,
    saveButtonProps = { labelId: "save", iconName: "Save" },
    localError,
    errorParseRules = [],
    fallbackErrorMessageId = "generic-error.unknown-error",
    farItems,
    messageBar: propsMessageBar,
  } = props;

  const { themes } = useSystemConfig();
  const { renderToString } = useContext(Context);

  const onFormSubmit = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      save();
    },
    [save]
  );

  const [isResetDialogVisible, setIsResetDialogVisible] = useState(false);
  const onDismissResetDialog = useCallback(() => {
    setIsResetDialogVisible(false);
  }, []);
  const doReset = useCallback(() => {
    reset();
    // If the form contains a CodeEditor, dialog dismiss animation does not play.
    // Defer the dismissal to ensure dismiss animation.
    setTimeout(() => setIsResetDialogVisible(false), 0);
  }, [reset]);

  const disabled = isUpdating || !isDirty;
  const commandBarItems: ICommandBarItemProps[] = useMemo(() => {
    return [
      {
        key: "save",
        text: renderToString(saveButtonProps.labelId),
        iconProps: { iconName: saveButtonProps.iconName },
        disabled: disabled || !canSave,
        onClick: () => save(),
      },
      {
        key: "reset",
        text: renderToString("reset"),
        iconProps: { iconName: "Delete" },
        disabled,
        theme: disabled ? themes.main : themes.destructive,
        onClick: () => setIsResetDialogVisible(true),
      },
    ];
  }, [canSave, disabled, save, saveButtonProps, renderToString, themes]);

  const resetDialogContentProps = useMemo(() => {
    return {
      title: <FormattedMessage id="FormContainer.reset-dialog.title" />,
      subText: renderToString("FormContainer.reset-dialog.message"),
    };
  }, [renderToString]);

  const {
    otherError,
    unhandledCauses: rawUnhandledCauses,
    value: { registerField, causes },
  } = useValidationError(updateError ?? localError);

  const [messageBar, formContext] = useMemo(() => {
    if (!otherError) {
      return [null, { registerField, causes }];
    }

    const {
      standaloneErrorMessageIds,
      fieldErrorMessageIds,
      unrecognizedError,
      unhandledCauses,
    } = parseError(otherError, rawUnhandledCauses, errorParseRules);

    const errorMessageIds = standaloneErrorMessageIds.slice();
    if (unrecognizedError) {
      errorMessageIds.push(fallbackErrorMessageId);
    }

    let messageBar: React.ReactNode;
    if (unhandledCauses.length > 0) {
      messageBar = (
        <ShowUnhandledValidationErrorCause causes={unhandledCauses} />
      );
    } else if (unrecognizedError) {
      messageBar = <ShowError error={unrecognizedError} />;
    } else if (errorMessageIds.length > 0) {
      messageBar = (
        <MessageBar messageBarType={MessageBarType.error}>
          {errorMessageIds.map((id, i) => (
            <FormattedMessage key={i} id={id} />
          ))}
        </MessageBar>
      );
    } else {
      messageBar = null;
    }

    const mappedCauses = { ...causes };
    for (const [field, messageIds] of Object.entries(fieldErrorMessageIds)) {
      const fieldCauses = mappedCauses[field] ?? [];
      for (const id of messageIds) {
        fieldCauses.push({
          kind: "general",
          location: field,
          details: { msg: renderToString(id) },
        });
      }
      mappedCauses[field] = fieldCauses;
    }

    return [messageBar, { registerField, causes: mappedCauses }];
  }, [
    otherError,
    rawUnhandledCauses,
    errorParseRules,
    causes,
    fallbackErrorMessageId,
    registerField,
    renderToString,
  ]);

  return (
    <FormContext.Provider value={formContext}>
      <CommandBarContainer
        isLoading={isUpdating}
        items={commandBarItems}
        farItems={farItems}
        messageBar={messageBar ?? propsMessageBar}
      >
        <form onSubmit={onFormSubmit}>{props.children}</form>
      </CommandBarContainer>
      <Dialog
        hidden={!isResetDialogVisible}
        dialogContentProps={resetDialogContentProps}
        onDismiss={onDismissResetDialog}
      >
        <DialogFooter>
          <PrimaryButton onClick={doReset} theme={themes.destructive}>
            <FormattedMessage id="reset" />
          </PrimaryButton>
          <DefaultButton onClick={onDismissResetDialog}>
            <FormattedMessage id="cancel" />
          </DefaultButton>
        </DialogFooter>
      </Dialog>
      <NavigationBlockerDialog blockNavigation={isDirty} />
    </FormContext.Provider>
  );
};

export default FormContainer;
