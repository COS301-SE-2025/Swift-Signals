describe('Simulations Page', () => {
  const API_BASE_URL = 'http://localhost:9090';
  const mockAuthToken = 'mock-auth-token-123';

  beforeEach(() => {
    cy.window().then((win) => {
      win.localStorage.setItem('authToken', mockAuthToken);
    });

    cy.intercept('GET', `${API_BASE_URL}/intersections`, {
      fixture: 'intersections.json'
    }).as('getIntersections');

    cy.intercept('POST', `${API_BASE_URL}/intersections`, {
      statusCode: 201,
      body: { id: 'new-intersection-123', message: 'Intersection created' }
    }).as('createIntersection');

    cy.intercept('GET', `${API_BASE_URL}/intersections/*/simulate`, {
      statusCode: 200,
      body: { message: 'Simulation started' }
    }).as('runSimulation');

    cy.visit('/simulations');
    cy.wait('@getIntersections');
  });

  it('renders the page with simulations table', () => {
    cy.contains('Recent Simulations').should('be.visible');
    cy.get('.simulationTable thead').within(() => {
      cy.contains('Name');
      cy.contains('Status');
      cy.contains('Actions');
    });
    cy.get('.simulationTable tbody tr').should('have.length.at.least', 1);
  });

  it('handles API error gracefully', () => {
    cy.intercept('GET', `${API_BASE_URL}/intersections`, {
      statusCode: 500,
      body: { error: 'Server error' }
    }).as('getError');

    cy.visit('/simulations');
    cy.wait('@getError');
    cy.contains('Error:').should('be.visible');
  });
  
});
