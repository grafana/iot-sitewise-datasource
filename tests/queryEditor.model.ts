import { type Page, type Locator } from '@playwright/test';

/**
 * Query Editor page object model testing utility providing properties and
 * methods for common locators and actions.
 */
export class QueryEditor {
  readonly queryTypeSelect: Locator;
  readonly assetSelect: Locator;
  readonly propertySelect: Locator;
  readonly propertyAliasInput: Locator;
  readonly qualitySelect: Locator;
  readonly formatSelect: Locator;
  readonly timeSelect: Locator;
  readonly aggregatePicker: Locator;
  readonly resolutionSelect: Locator;
  readonly modelIdSelect: Locator;
  readonly filterSelect: Locator;
  readonly showSelect: Locator;
  readonly queryOptionsSelect: Locator;
  readonly #page: Page;
  readonly featureToggles: Record<string, boolean>;

  constructor(page: Page, featureToggles: Record<string, boolean>) {
    this.queryTypeSelect = page.getByLabel('Query type');
    this.assetSelect = page.getByLabel('Asset');
    this.propertySelect = page.getByLabel('Property', { exact: true });
    this.propertyAliasInput = page.getByLabel('Property alias');
    this.qualitySelect = page.getByLabel('Quality');
    this.timeSelect = page.getByLabel('Time', { exact: true });
    this.formatSelect = page.getByLabel('Format');
    this.aggregatePicker = page.getByLabel('Aggregate');
    this.resolutionSelect = page.getByLabel('Resolution');
    this.modelIdSelect = page.getByLabel('Model ID');
    this.filterSelect = page.getByLabel('Filter', { exact: true });
    this.showSelect = page.getByLabel('Show', { exact: true });
    this.queryOptionsSelect = page.getByTestId('collapse-title');
    this.#page = page;
    this.featureToggles = featureToggles;
  }

  async selectQueryType(queryTypeOption: string) {
    // Open query type select options
    await this.queryTypeSelect.click();

    // Set query type
    await this.#page.getByText(queryTypeOption, { exact: true }).click();
  }

  async selectAsset(assetName: string) {
    await this.assetSelect.waitFor();
    await this.assetSelect.click();

    const assetOption = this.#page.getByText(assetName, { exact: true });

    await assetOption.waitFor();
    await assetOption.click();
  }

  async selectProperty(propertyName: string) {
    await this.propertySelect.waitFor();
    await this.propertySelect.click();

    const propertyOption = this.#page.getByText(propertyName, { exact: true });

    await propertyOption.waitFor();
    await propertyOption.click();
  }

  async selectShow(show: '** Parent **' | '** All **') {
    await this.showSelect.click();

    const showOption = this.#page.getByText(show, { exact: true });

    await showOption.click();
  }

  async runQuery() {
    const runQueryButton = this.#page.getByText('Run queries');
    await runQueryButton.click();
  }

  async openQueryOptions() {
    const optionsButton = this.#page.getByTestId('collapse-title');
    await optionsButton.click();
  }
}
