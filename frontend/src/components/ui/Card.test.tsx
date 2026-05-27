import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { Card } from './Card'

describe('Card', () => {
  it('renders children', () => {
    render(<Card><p>Content</p></Card>)
    expect(screen.getByText('Content')).toBeInTheDocument()
  })

  it('renders title when provided', () => {
    render(<Card title="My Title"><p>Content</p></Card>)
    expect(screen.getByText('My Title')).toBeInTheDocument()
  })

  it('does not render title when not provided', () => {
    const { container } = render(<Card><p>Content</p></Card>)
    const titles = container.querySelectorAll('h3')
    expect(titles.length).toBe(0)
  })

  it('renders subtitle when provided', () => {
    render(<Card title="Title" subtitle="Subtext"><p>Content</p></Card>)
    expect(screen.getByText('Subtext')).toBeInTheDocument()
  })
})
